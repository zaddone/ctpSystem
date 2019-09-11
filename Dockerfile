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
#RUN nohup ./main -db=/data/ins.db > /data/main.log 2>&1 &
#RUN nohup ./ctpServer 9999 150797 Dimon2019 tcp://218.202.237.33:10102 tcp://218.202.237.33:10112 >/data/ctp.log 2>&1 &
CMD ["sh","/data/binrun.sh"]
