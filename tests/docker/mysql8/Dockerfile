FROM centos:centos7

# Install dependencies
COPY install_dependencies_mysql8.sh /dist/install_dependencies_mysql8.sh
RUN /dist/install_dependencies_mysql8.sh
ENV PATH /usr/local/go/bin:/root/go/bin/:$PATH
