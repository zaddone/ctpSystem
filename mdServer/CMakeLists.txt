cmake_minimum_required(VERSION 2.8)
set( CMAKE_CXX_FLAGS "-std=c++11" )
project(mdServer)
set(SRC
${PROJECT_SOURCE_DIR}/main.cpp
${PROJECT_SOURCE_DIR}/marketspi.cpp
${PROJECT_SOURCE_DIR}/../api_tradeapi/socketunixserver.cpp
)
SET(HEADER
${PROJECT_SOURCE_DIR}/../api_tradeapi/ThostFtdcMdApi.h
${PROJECT_SOURCE_DIR}/../api_tradeapi/ThostFtdcUserApiDataType.h
${PROJECT_SOURCE_DIR}/../api_tradeapi/ThostFtdcUserApiStruct.h
${PROJECT_SOURCE_DIR}/marketspi.h
${PROJECT_SOURCE_DIR}/../api_tradeapi/socketunixserver.h
)
include_directories(
${PROJECT_SOURCE_DIR}
${PROJECT_SOURCE_DIR}/../api_tradeapi
)
LINK_DIRECTORIES(
${PROJECT_SOURCE_DIR}/../api_tradeapi
)
add_executable(${PROJECT_NAME} ${SRC} ${HEADER})
target_link_libraries(${PROJECT_NAME} -lthostmduser -lpthread)
