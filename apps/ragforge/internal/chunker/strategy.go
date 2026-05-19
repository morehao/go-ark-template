package chunker

import (
	"context"
	"strings"

	"github.com/morehao/golib/glog"
)

// Strategy values for SplitterConfig.Strategy.
const (
	StrategyAuto      = "auto"
	StrategyHeading   = "heading"
	StrategyHeuristic = "heuristic"
	StrategyRecursive = "recursive"
	StrategyLegacy    = "legacy"
)

// Split chunks text using the strategy configured in cfg.
func Split(text string, cfg SplitterConfig) []Chunk {
	if text == "" {
		return nil
	}
	cfg = ensureDefaults(cfg)
	chain, profile := resolveChainWithProfile(text, cfg)
	totalChars := len([]rune(text))

	var lastOut []Chunk
	for i, tier := range chain {
		out := runTier(tier, text, cfg, profile)
		if v := ValidateChunks(out, totalChars, cfg.ChunkSize); v.OK {
			return out
		}
		if tier == TierLegacy && i == len(chain)-1 {
			lastOut = out
		}
	}
	if lastOut != nil {
		return lastOut
	}
	return SplitText(text, cfg)
}

// TierRejection records why a tier was rejected by the validator.
type TierRejection struct {
	Tier   StrategyTier `json:"tier"`
	Reason string       `json:"reason"`
}

// Diagnostics captures which tier produced the returned chunks plus the
// chain that was attempted, any rejected tiers along the way, and the
// document profile that drove tier selection.
type Diagnostics struct {
	SelectedTier StrategyTier    `json:"selected_tier"`
	TierChain    []StrategyTier  `json:"tier_chain"`
	Rejected     []TierRejection `json:"rejected"`
	Profile      *DocProfile     `json:"profile,omitempty"`
}

// SplitWithDiagnostics is the same as Split but also returns the
// diagnostic trace (selected tier, full chain, rejection reasons, profile).
func SplitWithDiagnostics(text string, cfg SplitterConfig) ([]Chunk, *Diagnostics) {
	diag := &Diagnostics{SelectedTier: TierLegacy}
	if text == "" {
		return nil, diag
	}
	cfg = ensureDefaults(cfg)
	chain, profile := resolveChainWithProfile(text, cfg)
	diag.TierChain = chain
	diag.Profile = profile
	totalChars := len([]rune(text))

	var lastOut []Chunk
	var lastTier StrategyTier
	for i, tier := range chain {
		out := runTier(tier, text, cfg, profile)
		v := ValidateChunks(out, totalChars, cfg.ChunkSize)
		if v.OK {
			diag.SelectedTier = tier
			return out, diag
		}
		diag.Rejected = append(diag.Rejected, TierRejection{Tier: tier, Reason: v.Reason})
		glog.Infof(context.Background(), "chunker: tier %s rejected: %s", tier, v.Reason)
		if tier == TierLegacy && i == len(chain)-1 {
			lastOut = out
			lastTier = tier
		}
	}
	if lastOut != nil {
		diag.SelectedTier = lastTier
		return lastOut, diag
	}
	return SplitText(text, cfg), diag
}

// SplitParentChild is the strategy-aware analog of SplitTextParentChild.
func SplitParentChild(text string, parentCfg, childCfg SplitterConfig) ParentChildResult {
	if text == "" {
		return ParentChildResult{}
	}
	parentCfg = ensureDefaults(parentCfg)
	childCfg = ensureDefaults(childCfg)

	parents := Split(text, parentCfg)
	if len(parents) == 0 {
		return ParentChildResult{}
	}

	var newParents []Chunk
	var children []ChildChunk
	childSeq := 0
	for _, parent := range parents {
		subs := Split(parent.Content, childCfg)

		parentIndex := -1
		if len(subs) > 1 || (len(subs) == 1 && subs[0].Content != parent.Content) {
			parentIndex = len(newParents)
			newParents = append(newParents, parent)
		}
		for _, sub := range subs {
			sub.Seq = childSeq
			sub.Start += parent.Start
			sub.End += parent.Start
			sub.ContextHeader = mergeBreadcrumbs(parent.ContextHeader, sub.ContextHeader)
			children = append(children, ChildChunk{Chunk: sub, ParentIndex: parentIndex})
			childSeq++
		}
	}
	return ParentChildResult{Parents: newParents, Children: children}
}

func mergeBreadcrumbs(parent, child string) string {
	if parent == "" {
		return child
	}
	if child == "" {
		return parent
	}
	parentLines := strings.Split(parent, "\n")
	childLines := strings.Split(child, "\n")
	if len(parentLines) > 0 && len(childLines) > 0 &&
		strings.TrimSpace(parentLines[len(parentLines)-1]) == strings.TrimSpace(childLines[0]) {
		childLines = childLines[1:]
	}
	if len(childLines) == 0 {
		return parent
	}
	return parent + "\n" + strings.Join(childLines, "\n")
}

func resolveChainWithProfile(text string, cfg SplitterConfig) ([]StrategyTier, *DocProfile) {
	switch cfg.Strategy {
	case StrategyHeading:
		return []StrategyTier{TierHeading, TierLegacy}, nil
	case StrategyHeuristic:
		return []StrategyTier{TierHeuristic, TierLegacy}, nil
	case StrategyRecursive:
		return []StrategyTier{TierLegacy}, nil
	case StrategyLegacy, "":
		return []StrategyTier{TierLegacy}, nil
	case StrategyAuto:
		fallthrough
	default:
		profile := ProfileDocument(text)
		return SelectStrategy(profile), profile
	}
}

func runTier(tier StrategyTier, text string, cfg SplitterConfig, profile *DocProfile) []Chunk {
	switch tier {
	case TierHeading:
		return splitByHeadings(text, cfg, profile)
	case TierHeuristic:
		return splitByHeuristics(text, cfg, profile)
	case TierLegacy:
		return SplitText(text, cfg)
	}
	return SplitText(text, cfg)
}

func ensureDefaults(cfg SplitterConfig) SplitterConfig {
	if cfg.ChunkSize <= 0 {
		cfg.ChunkSize = DefaultChunkSize
	}
	if cfg.ChunkOverlap <= 0 {
		cfg.ChunkOverlap = DefaultChunkOverlap
	}
	if len(cfg.Separators) == 0 {
		cfg.Separators = []string{"\n\n", "\n", "。"}
	}
	if cfg.TokenLimit > 0 {
		lang := LangMixed
		if len(cfg.Languages) > 0 {
			lang = cfg.Languages[0]
		}
		charBudget := CharsForTokenLimit(cfg.TokenLimit, lang)
		if charBudget > 0 && (cfg.ChunkSize == 0 || charBudget < cfg.ChunkSize) {
			cfg.ChunkSize = charBudget
		}
	}
	if cfg.ChunkOverlap > cfg.ChunkSize/2 && cfg.ChunkSize > 0 {
		cfg.ChunkOverlap = cfg.ChunkSize / 2
	}
	return cfg
}

var splitByHeadings = func(text string, cfg SplitterConfig, _ *DocProfile) []Chunk {
	return SplitText(text, cfg)
}

var splitByHeuristics = func(text string, cfg SplitterConfig, _ *DocProfile) []Chunk {
	return SplitText(text, cfg)
}
