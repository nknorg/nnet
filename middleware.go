package nnet

import "github.com/nknorg/nnet/util"

// ApplyMiddleware add a middleware to overlay or local node
func (nn *NNet) ApplyMiddleware(f interface{}) error {
	applied := false
	errs := util.NewErrors()

	err := nn.overlay.ApplyMiddleware(f)
	if err == nil {
		applied = true
	} else {
		errs = append(errs, err)
	}

	err = nn.localNode.ApplyMiddleware(f)
	if err == nil {
		applied = true
	} else {
		errs = append(errs, err)
	}

	if !applied {
		return errs.Merged()
	}

	return nil
}
