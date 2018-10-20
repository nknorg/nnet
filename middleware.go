package nnet

import "github.com/nknorg/nnet/util"

// ApplyMiddleware add a middleware to overlay or local node
func (nn *NNet) ApplyMiddleware(f interface{}) error {
	applied := false
	errs := util.NewErrors()

	err := nn.LocalNode.ApplyMiddleware(f)
	if err == nil {
		applied = true
	} else {
		errs = append(errs, err)
	}

	err = nn.Overlay.ApplyMiddleware(f)
	if err == nil {
		applied = true
	} else {
		errs = append(errs, err)
	}

	for _, router := range nn.Overlay.GetRouters() {
		err = router.ApplyMiddleware(f)
		if err == nil {
			applied = true
		} else {
			errs = append(errs, err)
		}
	}

	if !applied {
		return errs.Merged()
	}

	return nil
}
