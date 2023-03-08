.PHONY: deploy
deploy:
	pushd infra && \
		aws-vault exec default -- cdk deploy --require-approval never && \
		popd
