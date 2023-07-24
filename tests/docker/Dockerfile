FROM centos:centos7

# Install dependencies
COPY install_dependencies.sh /dist/install_dependencies.sh
RUN /dist/install_dependencies.sh
ENV PATH /usr/local/go/bin:/root/go/bin:$PATH
