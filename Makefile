default: test

test:
	ginkgo -r -v -cover
	# go test $$(glide novendor)

fmt:
	 go fmt $$(glide novendor)
cover:
	ginkgo -cover -r -v
