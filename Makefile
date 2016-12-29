default: test

test:
	ginkgo -r -v
	# go test $$(glide novendor)

cover:
	ginkgo -cover -r -v
