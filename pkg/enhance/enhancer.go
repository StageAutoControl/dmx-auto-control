package enhance

import "github.com/StageAutoControl/controller/pkg/cntl"

// Enhancers stores the globally registered enhancers
var Enhancers = make([]cntl.Enhancer, 0)