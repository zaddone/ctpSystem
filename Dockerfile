#
# ctp server
#
# https://github.com/zaddone/ctpSystem
#

# Pull base image.
FROM ubuntu:latest
MAINTAINER zaddone@qq.com
ADD . /code
RUN echo /code >> /etc/ld.so.conf
RUN ldconfig
RUN mkdir /data
WORKDIR /code
CMD ["sh","/data/binrun.sh"]
