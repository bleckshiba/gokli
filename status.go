package main

type Status int

const (
	PENDING Status = iota
	WIP
	DONE
)

func (status Status) getNext() Status {
	if status == DONE {
		return PENDING
	}
	return status + 1
}

func (status Status) getPrev() Status {
	if status == PENDING {
		return DONE
	}
	return status - 1
}
