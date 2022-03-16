STACK_NAME := sam-app-stack
STACK_BUCKET := mysd33bucket123sam

.PHONY: build
.PHONY: unit_test
.PHONY: deploy
.PHONY: deploy_guided
.PHONY: delete

build:
# sam build
# for windows
	sam.cmd build
	xcopy /I config .aws-sam\build\GetUsersFunction\config
	xcopy /I config .aws-sam\build\PostUsersFunction\config

unit_test:
	cd ap & go test -v ./internal/...

deploy_guided:
# for windows
	sam.cmd deploy --guided

deploy:
# for windows
	sam.cmd deploy

delete:
	aws cloudformation delete-stack --stack-name $(STACK_NAME)
	aws s3 rm "s3://$(STACK_BUCKET)" --recursive