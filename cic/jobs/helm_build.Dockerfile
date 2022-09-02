FROM alpine/git AS CLONE_SOURCE
RUN mkdir /stuff
WORKDIR /stuff
RUN git clone https://$gitRegistry/$gitRepoName
WORKDIR /stuff/$gitRepoName

FROM helm AS BUILD_HELM_PACKAGE
RUN mkdir /stuff
WORKDIR /stuff
COPY --from=CLONE_SOURCE /stuff/$gitRepoName /stuff
RUN  helm package ./install --version $helmPackageVersion
RUN  curl --data-binary "@${FN}" https://$helmRegistry/api/charts