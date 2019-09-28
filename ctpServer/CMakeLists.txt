cmake_minimum_required(VERSION 2.8)
set( CMAKE_CXX_FLAGS "-std=c++11" )
project(ctpServer)
#AUX_SOURCE_DIRECTORY(${PROJECT_SOURCE_DIR} DIR_SRCS)
include_directories(
${PROJECT_SOURCE_DIR}
${PROJECT_SOURCE_DIR}/api_tradeapi
)
set(SRC
${PROJECT_SOURCE_DIR}/main.cpp
${PROJECT_SOURCE_DIR}/ctpspi.cpp
${PROJECT_SOURCE_DIR}/marketspi.cpp
${PROJECT_SOURCE_DIR}/socketunixserver.cpp
${PROJECT_SOURCE_DIR}/traderspi.cpp
)
SET(HEADER
${PROJECT_SOURCE_DIR}/api_tradeapi/ThostFtdcMdApi.h
${PROJECT_SOURCE_DIR}/api_tradeapi/ThostFtdcTraderApi.h
${PROJECT_SOURCE_DIR}/api_tradeapi/ThostFtdcUserApiDataType.h
${PROJECT_SOURCE_DIR}/api_tradeapi/ThostFtdcUserApiStruct.h
${PROJECT_SOURCE_DIR}/ctpspi.h
${PROJECT_SOURCE_DIR}/marketspi.h
${PROJECT_SOURCE_DIR}/socketunixserver.h
${PROJECT_SOURCE_DIR}/traderspi.h
)
LINK_DIRECTORIES(${PROJECT_SOURCE_DIR}/api_tradeapi)
add_executable(${PROJECT_NAME} ${SRC} ${HEADER})
target_link_libraries(${PROJECT_NAME} -lthostmduser -lthosttraderapi -lpthread)
