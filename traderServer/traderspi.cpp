#include <iostream>
#include <unistd.h>
//#include <sys/stat.h>
#include <thread>
//#include <string.h>
#include "traderspi.h"

using namespace std;
bool IsFlowControl(int iResult)
{
    return ((iResult == -2) || (iResult == -3));
}

TraderSpi::TraderSpi(const char * path):socketUnixServer(path){
    //if (0 != access(path,0)){
    //    mkdir(path,0777);
    //}
    this->path = path;

    memset(&this->userReq,0,sizeof(this->userReq));
    this->initMap();
}
TraderSpi::TraderSpi(
        const char * brokerID,
        const char * userID,
        const char *password,
        const char *passwordBak,
        const char *addr,
        const char * path):socketUnixServer(path){

    this->trApi = NULL;
    this->TradingDay = NULL;
    this->queryIns = false;
    this->path = path;
    this->Addr = addr;
    this->initMap();
    memset(&this->userReq,0,sizeof(this->userReq));
    this->setUserReg(brokerID,userID,password,passwordBak);
    this->run();
}

void TraderSpi::initMap(){
    //this->mapstring["ins"] = 1;
    this->mapstring["stop"] = 100;
    this->mapstring["config"] = 2;
    this->mapstring["addr"] = 3;
    this->mapstring["ReqQrySettlementInfo"] = 4;
    this->mapstring["ReqSettlementInfoConfirm"] = 5;
    this->mapstring["ReqQrySettlementInfoConfirm"] = 6;
    this->mapstring["TradingAccount"] = 7;
    this->mapstring["InvestorPositionDetail"] = 8;
    this->mapstring["InvestorPosition"] = 9;
    this->mapstring["open"] = 10;
    this->mapstring["close"] = 11;
}

void TraderSpi::routeHand(const char *data){

    //if (NULL==this->TradingDay){
    //    cout<<"trader:"<<data<<endl;
    //    return;
    //}
    char db[1024];
    strcpy(db,data);

    //cout<<"db "<<db<<endl;
    char *p;
    char sep[] = " ";
    char str[100][1024];
    p = strtok(db,sep);
    int i;
    i = 0;
    while( p != NULL ) {
        strcpy(str[i] , p);
        p = strtok(NULL, sep);
        i++;
    }
    switch (this->mapstring[str[0]]){
    case 100:{
        this->stop();
        cout<<"stop"<<endl;
    }
        break;
    case 2:{
        this->setUserReg(str[1],str[2],str[3],str[4]);
        this->Addr = str[5];
        this->run();
    }
        break;
    case 3:{
        //this->setUserReg(str[1],str[2],str[3],str[4]);
        this->Addr = str[1];
        this->run();
    }
        break;
    case 4:{

        this->reqQrySettlementInfo();
    }
        break;
    case 5:{
        this->reqSettlementInfoConfirm();
    }
        break;
    case 6:{
        this->reqQrySettlementInfoConfirm();
    }
        break;
    case 7:{
        cout<<"db:"<<db<<endl;
        this->reqTradingAccount();
    }
        break;
    case 8:{
        char * ins=NULL;
        if (i>1){
            ins = str[1];
        }
        this->reqInvestorPositionDetail(ins);
    }
        break;
    case 9:{
        char * ins=NULL;
        if (i>1){
            ins = str[1];
        }
        this->reqInvestorPosition(ins);
    }
        break;

    case 10:{
        char *dis;
        if (i>1)dis= str[2];
        double pr=0;
        if (i>2){
            pr = atof(str[3]);
        }
        this->sendOrderOpen(str[1],dis,pr);
    }
        break;
    case 11:{
        this->sendOrderClose(str[1]);
    }
        break;
    default:
        printf("default %s %s end",data,str[0]);
        break;
    }
}


void TraderSpi::OnFrontDisconnected(int nReason){

    cout<<"disconnected:"<<nReason<<endl;
    this->stop();
}
void TraderSpi::OnFrontConnected(){
    cout<<"conn"<<endl;
    this->reqUserLogin();

}
void TraderSpi::OnRspQryInvestorPosition(
            CThostFtdcInvestorPositionField *pInvestorPosition,
            CThostFtdcRspInfoField *pRspInfo,
            int nRequestID,
        bool bIsLast){
    if (pRspInfo && pRspInfo->ErrorID!=0){
        cout<<pRspInfo->ErrorMsg<<endl;
        return;
    }
    if (!bIsLast)return;
    if (!pInvestorPosition)return;
    //this->send(pInvestorPosition->InstrumentID);
    cout<<"InstrumentID:"<<pInvestorPosition->InstrumentID<<endl;
    cout<<"PositionDate:"<<pInvestorPosition->PositionDate<<endl;
    cout<<"LongFrozenAmount:"<<pInvestorPosition->LongFrozenAmount<<endl;
    cout<<"ShortFrozenAmount:"<<pInvestorPosition->ShortFrozenAmount<<endl;
    cout<<"OpenAmount:"<<pInvestorPosition->OpenAmount<<endl;
    cout<<"CloseAmount:"<<pInvestorPosition->CloseAmount<<endl;
    cout<<"PositionCost:"<<pInvestorPosition->PositionCost<<endl;
    cout<<"UseMargin:"<<pInvestorPosition->UseMargin<<endl;
    cout<<"Commission:"<<pInvestorPosition->Commission<<endl;
    cout<<"CloseProfit:"<<pInvestorPosition->CloseProfit<<endl;
    cout<<"PositionProfit:"<<pInvestorPosition->PositionProfit<<endl;
    cout<<"TradingDay:"<<pInvestorPosition->TradingDay<<endl;


}

void TraderSpi::OnRspQryTradingAccount(
        CThostFtdcTradingAccountField *pTradingAccount,
        CThostFtdcRspInfoField *pRspInfo,
        int nRequestID,
        bool bIsLast) {
    //return;
    cout<<"onRspTA"<<endl;
    if (pRspInfo && pRspInfo->ErrorID!=0){
        cout<<pRspInfo->ErrorMsg<<endl;
        return;
    }

    if (!bIsLast)return;

    if (!pTradingAccount)return;
    cout<<pTradingAccount->Deposit<<endl;

}



void TraderSpi::OnRspQryInvestorPositionDetail(
        CThostFtdcInvestorPositionDetailField *pInvestorPositionDetail,
        CThostFtdcRspInfoField *pRspInfo,
        int nRequestID,
        bool bIsLast){
    if (pRspInfo && pRspInfo->ErrorID!=0){
        cout<<pRspInfo->ErrorMsg<<endl;
        return;
    }
    cout <<"InstrumentID "<< pInvestorPositionDetail->InstrumentID << endl;
    cout <<"HedgeFlag " << pInvestorPositionDetail->HedgeFlag << endl;
    cout <<"Direction " << pInvestorPositionDetail->Direction << endl;
    cout <<"OpenDate " << pInvestorPositionDetail->OpenDate << endl;
    cout <<"TradeID " << pInvestorPositionDetail->TradeID << endl;
    cout <<"Volume " << pInvestorPositionDetail->Volume << endl;
    cout <<"OpenPrice " << pInvestorPositionDetail->OpenPrice << endl;
    cout <<"TradingDay " << pInvestorPositionDetail->TradingDay << endl;
    cout <<"SettlementID " << pInvestorPositionDetail->SettlementID << endl;
    cout <<"TradeType " << pInvestorPositionDetail->TradeType << endl;
    cout <<"CombInstrumentID " << pInvestorPositionDetail->CombInstrumentID << endl;
    cout <<"ExchangeID " << pInvestorPositionDetail->ExchangeID << endl;
    cout <<"Margin " << pInvestorPositionDetail->Margin << endl;
    cout <<"ExchMargin " << pInvestorPositionDetail->ExchMargin << endl;
    cout <<"MarginRateByMoney " << pInvestorPositionDetail->MarginRateByMoney << endl;
    cout <<"MarginRateByVolume " << pInvestorPositionDetail->MarginRateByVolume << endl;
    cout <<"LastSettlementPrice " << pInvestorPositionDetail->LastSettlementPrice << endl;
    cout <<"SettlementPrice " << pInvestorPositionDetail->SettlementPrice << endl;
    cout <<"CloseVolume " << pInvestorPositionDetail->CloseVolume << endl;
    cout <<"CloseAmount " << pInvestorPositionDetail->CloseAmount << endl;


}

void TraderSpi::OnErrRtnOrderInsert(CThostFtdcInputOrderField *pInputOrder, CThostFtdcRspInfoField *pRspInfo) {

    if (pRspInfo && 0!=pRspInfo->ErrorID){
        cout<<pRspInfo->ErrorMsg<<endl;
        return;
    }
    cout<<"order err"<<pInputOrder->OrderRef<<endl;

}
void TraderSpi::OnRspOrderInsert(
        CThostFtdcInputOrderField *pInputOrder,
        CThostFtdcRspInfoField *pRspInfo,
        int nRequestID,
        bool bIsLast){
    if (pRspInfo && 0!=pRspInfo->ErrorID){
        cout<<pRspInfo->ErrorMsg<<endl;
        return;
    }
    cout<<"order insert"<<pInputOrder->OrderRef<<endl;
}

void TraderSpi::stop(){
    if (this->trApi== NULL) return;
    this->trApi->RegisterSpi(NULL);
    this->trApi->Release();
    //this->trApi->Join();
    //this->Join();
    this->trApi = NULL;
    this->over =  true;
    cout<<"stop ok"<<endl;
}
void TraderSpi::Join(){
    this->trApi->Join();
    //this->stop();
    cout<<"stop trader"<<endl;
    //this->send("addr");
}
int TraderSpi::getRequestID(){
    this->requestID++;
    return this->requestID;
}



void TraderSpi::OnRspQryInstrument(
        CThostFtdcInstrumentField *pInstrument,
        CThostFtdcRspInfoField *pRspInfo,
        int nRequestID,
        bool bIsLast)
{
    //cout<< pInstrument->InstrumentName<<endl;
    //cout<< pInstrument->InstrumentID<<endl;
    this->mapInstrument[pInstrument->InstrumentID] = *pInstrument;
    char db[100] = "ins ";
    strcat(db,pInstrument->InstrumentID);
    //cout<< "ins "<<pInstrument->InstrumentID <<endl;
    this->send(db);
    //cout<< db <<endl;
    //this->routeHand(db);
    //pInstrument->InstrumentID;
    //collect(pInstrument);

    //if (bIsLast)
        //signal(allInstrumentsReady);
}
void TraderSpi::swapPassword(){
    TThostFtdcBrokerIDType bakPass;
    strcpy(bakPass,this->userReq.Password);
    strcpy(this->userReq.Password,this->pass);
    strcpy(this->pass,bakPass);

}
void TraderSpi::OnRspUserLogin(
    CThostFtdcRspUserLoginField *pRspUserLogin,
    CThostFtdcRspInfoField *pRspInfo,
    int nRequestID,
    bool bIsLast)
{

    cout<<"trader"<<pRspInfo->ErrorID<<endl;

    //if (pRspInfo->ErrorID != 0 ){
    //    cout<<"trader:"<<pRspInfo->ErrorID<<endl;
    //    cout<<pRspInfo->ErrorMsg<<endl;

    //}
    //char pass[]="abc2019";
    if (140==pRspInfo->ErrorID){
        CThostFtdcUserPasswordUpdateField res;
        memset(&res,0,sizeof(res));
        strcpy(res.BrokerID,this->userReq.BrokerID);
        strcpy(res.UserID,this->userReq.UserID);
        strcpy(res.OldPassword,this->userReq.Password);
        strcpy(res.NewPassword,pass);
        if (0==this->trApi->ReqUserPasswordUpdate(&res,this->getRequestID())){
            this->swapPassword();
            //TThostFtdcBrokerIDType bakPass;
            //strcpy(bakPass,this->userReq.Password);
            //strcpy(this->userReq.Password,pass);
            //strcpy(pass,bakPass);
            this->reqUserLogin();

        }
    }else if (3 == pRspInfo->ErrorID){
        //strcpy(this->userReq.Password,pass);
        this->swapPassword();
        this->reqUserLogin();
    }else if (0 == pRspInfo->ErrorID){
        this->TradingDay = 	this->trApi->GetTradingDay();
        cout <<"Td connected "<<this->TradingDay << endl;
        this->reqInstruments();
    }else if (7 == pRspInfo->ErrorID){
        this->swapPassword();
        //TThostFtdcBrokerIDType bakPass;
        //strcpy(bakPass,this->userReq.Password);
        //strcpy(this->userReq.Password,pass);
        //strcpy(pass,bakPass);
        this->reqUserLogin();
    }else{
        cout<<pRspInfo->ErrorMsg<<endl;
    }
    //if (0 == pRspInfo->ErrorID){
    //    this->queryInstruments();
    //};
}

void TraderSpi::run(const char * addr){
    this->trApi = CThostFtdcTraderApi::CreateFtdcTraderApi(this->path);
    char _addr[1024];
    strcpy(_addr,addr);
    cout<<_addr<<endl;
    this->trApi->SubscribePublicTopic(THOST_TERT_RESTART);				// 注册公有流
    this->trApi->SubscribePrivateTopic(THOST_TERT_RESTART);
    this->trApi->RegisterFront(_addr);

    this->trApi->RegisterSpi(this);
    this->trApi->Init();
    //this->Join();
    thread th(&TraderSpi::Join,this);
    th.detach();

}

void TraderSpi::run(){
    if (this->trApi != NULL) return;
    this->trApi = CThostFtdcTraderApi::CreateFtdcTraderApi(this->path);
    this->trApi->RegisterSpi(this);
    char _addr[1024];
    strcpy(_addr,Addr);
    //cout<<_addr<<endl;
    this->trApi->SubscribePublicTopic(THOST_TERT_RESTART);				// 注册公有流
    this->trApi->SubscribePrivateTopic(THOST_TERT_RESTART);
    this->trApi->RegisterFront(_addr);
    this->trApi->Init();
    //this->Join();
    thread th(&TraderSpi::Join,this);
    th.detach();
    //mSpi->mdApi->Join();
}

void TraderSpi::setUserReg(
        const char * brokerID,
        const char * userID,
        const char *password,
        const char *passwordBak){
    //memset(&this->userReq,0,sizeof(this->userReq));
    strcpy(this->userReq.BrokerID,brokerID);
    strcpy(this->userReq.UserID,userID);
    strcpy(this->userReq.Password,password);
    strcpy(this->pass,passwordBak);
}

void TraderSpi::OnRspSettlementInfoConfirm(
    CThostFtdcSettlementInfoConfirmField *pSettlementInfoConfirm,
    CThostFtdcRspInfoField *pRspInfo,
    int nRequestID,
    bool bIsLast) {
    if (pRspInfo && 0!=pRspInfo->ErrorID){
        cout<<pRspInfo->ErrorMsg<<endl;
        return;
    }
    if (!bIsLast)return;
    if (!pSettlementInfoConfirm) return;
    cout<<pSettlementInfoConfirm->ConfirmDate << endl;
    cout<<pSettlementInfoConfirm->ConfirmTime << endl;
    cout<<pSettlementInfoConfirm->SettlementID << endl;
    //this->send(msg);
    return;


}

void TraderSpi::OnRspQrySettlementInfo(
    CThostFtdcSettlementInfoField *pSettlementInfo,
    CThostFtdcRspInfoField *pRspInfo,
    int nRequestID,
    bool bIsLast){
    if (pRspInfo && 0!=pRspInfo->ErrorID){
        cout<<pRspInfo->ErrorMsg<<endl;
        return;
    }

    //cout<<"bIsLast"<<bIsLast<<endl;
    if (!bIsLast)return;
    //char msg[8192];


    if (!pSettlementInfo) return;
    cout<<pSettlementInfo->TradingDay<<endl;
    cout<<pSettlementInfo->SequenceNo<< endl;
    cout<<pSettlementInfo->SettlementID<< endl;
    cout<<pSettlementInfo->Content<< endl;
    //sprintf(msg,
    //        "msg TradingDay|%s SequenceNo|%d SettlementID|%d Content|%s",
    //        pSettlementInfo->TradingDay,
    //        pSettlementInfo->SequenceNo,
    //        pSettlementInfo->SettlementID,
    //        pSettlementInfo->Content);
    //cout<<msg<<endl;
    //this->send(msg);
    return;


}

void TraderSpi::OnRspQrySettlementInfoConfirm(
        CThostFtdcSettlementInfoConfirmField *pSettlementInfoConfirm,
        CThostFtdcRspInfoField *pRspInfo,
        int nRequestID,
        bool bIsLast){

    if (pRspInfo && 0!=pRspInfo->ErrorID){
        cout<<pRspInfo->ErrorMsg<<endl;
        return;
    }
    if (!bIsLast)return;
    //char msg[8192];
    if (!pSettlementInfoConfirm)return;
    cout << pSettlementInfoConfirm-> ConfirmDate << endl;
    cout << pSettlementInfoConfirm-> ConfirmTime << endl;
    cout << pSettlementInfoConfirm-> SettlementID << endl;
    //this->send(msg);
    return;
}

void TraderSpi::reqInstruments()
{
    if (NULL==this->TradingDay)return;
    if (this->queryIns)return;
    this->queryIns = true;
    //cout<<"query ins"<<endl;
    CThostFtdcQryInstrumentField req;
    memset(&req, 0, sizeof(req));
    //this->trApi->ReqQryInstrument(&req,this->getRequestID());
    while (true)
    {
        //cout<<"query ins"<<endl;
        int iResult = this->trApi->ReqQryInstrument(&req, this->getRequestID());
        if (!IsFlowControl(iResult))
        {
            break;
        }
        else
        {
            sleep(1);
        }
    }
}

void TraderSpi::reqSettlementInfoConfirm(){
    if (NULL==this->TradingDay)return;
    CThostFtdcSettlementInfoConfirmField pSettlementInfoConfirm;
    memset(&pSettlementInfoConfirm,0,sizeof(pSettlementInfoConfirm));
    strcpy(pSettlementInfoConfirm.BrokerID,this->userReq.BrokerID);
    strcpy(pSettlementInfoConfirm.AccountID,this->userReq.UserID);
    strcpy(pSettlementInfoConfirm.InvestorID,this->userReq.UserID);

    while (true)
    {
        int iResult = this->trApi->ReqSettlementInfoConfirm(&pSettlementInfoConfirm,this->getRequestID());
        if (!IsFlowControl(iResult))
        {
            break;
        }
        else
        {
            sleep(1);
        }
    }
}
void TraderSpi::reqQrySettlementInfo(){
    if (NULL==this->TradingDay)return;
    CThostFtdcQrySettlementInfoField pQrySettlementInfo;
    memset(&pQrySettlementInfo,0,sizeof(pQrySettlementInfo));
    strcpy(pQrySettlementInfo.BrokerID,this->userReq.BrokerID);
    strcpy(pQrySettlementInfo.AccountID,this->userReq.UserID);
    strcpy(pQrySettlementInfo.InvestorID,this->userReq.UserID);
    while (true)
    {
        int iResult = this->trApi->ReqQrySettlementInfo(&pQrySettlementInfo,this->getRequestID());
        if (!IsFlowControl(iResult))
        {
            break;
        }
        else
        {
            sleep(1);
        }
    }
}
void TraderSpi::reqQrySettlementInfoConfirm(){

    if (NULL==this->TradingDay)return;
    CThostFtdcQrySettlementInfoConfirmField pQrySettlementInfoConfirm;
    memset(&pQrySettlementInfoConfirm,0,sizeof(pQrySettlementInfoConfirm));
    strcpy(pQrySettlementInfoConfirm.BrokerID,this->userReq.BrokerID);
    strcpy(pQrySettlementInfoConfirm.AccountID,this->userReq.UserID);
    strcpy(pQrySettlementInfoConfirm.InvestorID,this->userReq.UserID);
    while (true)
    {
        int iResult = this->trApi->ReqQrySettlementInfoConfirm(&pQrySettlementInfoConfirm,this->getRequestID());
        if (!IsFlowControl(iResult))
        {
            break;
        }
        else
        {
            sleep(1);
        }
    }
}
void TraderSpi::reqTradingAccount(){

    cout <<"Td connected "<<this->TradingDay << endl;
    if (NULL==this->TradingDay)return;
    CThostFtdcQryTradingAccountField pQryTradingAccount;
    memset(&pQryTradingAccount,0,sizeof(pQryTradingAccount));
    strcpy(pQryTradingAccount.AccountID,this->userReq.UserID);
    strcpy(pQryTradingAccount.BrokerID,this->userReq.BrokerID);
    strcpy(pQryTradingAccount.InvestorID,this->userReq.UserID);
    pQryTradingAccount.BizType=THOST_FTDC_BZTP_Future;
    while (true)
    {
        int iResult = this->trApi->ReqQryTradingAccount(&pQryTradingAccount,this->getRequestID());
        cout<<iResult<<endl;
        if (!IsFlowControl(iResult))
        {
            break;
        }
        else
        {
            sleep(1);
        }
    }

}

void TraderSpi::reqInvestorPosition(const char * ins){
    //return;
    if (NULL==this->TradingDay)return;
    CThostFtdcQryInvestorPositionField pQryInvestorPosition;
    memset(&pQryInvestorPosition,0,sizeof(pQryInvestorPosition));
    strcpy(pQryInvestorPosition.InvestorID,this->userReq.UserID);
    strcpy(pQryInvestorPosition.BrokerID,this->userReq.BrokerID);
    if (ins!=NULL) strcpy(pQryInvestorPosition.InstrumentID,ins);
    while (true)
    {
        int iResult = this->trApi->ReqQryInvestorPosition(&pQryInvestorPosition,this->getRequestID());
        cout<<"InvestorPosition "<<iResult<<endl;
        if (!IsFlowControl(iResult))
        {
            break;
        }
        else
        {
            sleep(1);
        }
    }
}

void TraderSpi::reqInvestorPositionDetail(const char * ins){

    if (NULL==this->TradingDay)return;
    CThostFtdcQryInvestorPositionDetailField pInvestorPositionDetail;
    memset(&pInvestorPositionDetail,0,sizeof(pInvestorPositionDetail));
    strcpy(pInvestorPositionDetail.BrokerID,this->userReq.BrokerID);
    strcpy(pInvestorPositionDetail.InvestorID,this->userReq.UserID);
    if (NULL!=ins)strcpy(pInvestorPositionDetail.InstrumentID,ins);

    while (true)
    {
        int iResult = this->trApi->ReqQryInvestorPositionDetail(&pInvestorPositionDetail,this->getRequestID());
        if (!IsFlowControl(iResult))
        {
            break;
        }
        else
        {
            sleep(1);
        }
    }

}

void TraderSpi::sendOrderClose(const char * ins){

    if (NULL==this->TradingDay)return;
    CThostFtdcInstrumentField insinfo = this->mapInstrument[ins];
    cout<<"close "<<ins<<endl;
    CThostFtdcInputOrderField order;
    memset(&order,0,sizeof(order));
    strcpy(order.BrokerID,this->userReq.BrokerID);
    strcpy(order.InvestorID,this->userReq.UserID);
    strcpy(order.InstrumentID,ins);
    //strcpy(order.OrderRef,"test1");
    strcpy(order.UserID,this->userReq.UserID);
    strcpy(order.ExchangeID,insinfo.ExchangeID);
    order.ContingentCondition =THOST_FTDC_CC_Immediately;

    order.CombOffsetFlag[0] = THOST_FTDC_OF_CloseToday;
    order.CombOffsetFlag[1] = THOST_FTDC_OF_CloseYesterday;

    order.CombHedgeFlag[0] = THOST_FTDC_HF_Speculation;

    order.VolumeTotalOriginal = 1;
    order.VolumeCondition = THOST_FTDC_VC_AV;
    order.MinVolume = 1;
    order.ForceCloseReason = THOST_FTDC_FCC_NotForceClose;
    order.IsAutoSuspend = 0;
    order.UserForceClose = 0;

    while (true)
    {
        int iResult = this->trApi->ReqOrderInsert(&order,this->getRequestID());
        if (!IsFlowControl(iResult))
        {
            break;
        }
        else
        {
            sleep(1000);
        }
    }

}

void TraderSpi::sendOrderOpen(const char *ins, const char *dir,const double price){

    if (NULL==this->TradingDay)return;
    CThostFtdcInstrumentField insinfo = this->mapInstrument[ins];
    //cout<<insinfo.InstrumentName<<endl;
    cout<<"open "<<ins<<endl;
    CThostFtdcInputOrderField order;
    memset(&order,0,sizeof(order));
    strcpy(order.BrokerID,this->userReq.BrokerID);
    strcpy(order.InvestorID,this->userReq.UserID);
    strcpy(order.InstrumentID,ins);
    //strcpy(order.OrderRef,"test1");
    strcpy(order.UserID,this->userReq.UserID);
    strcpy(order.ExchangeID,insinfo.ExchangeID);
    order.ContingentCondition =THOST_FTDC_CC_Immediately;
    //order.StopPrice=
    if (dir=="buy"){
        order.Direction = THOST_FTDC_D_Buy;
    }else{
        order.Direction = THOST_FTDC_D_Sell;
    }
    order.CombOffsetFlag[0] = THOST_FTDC_OF_Open;

    order.CombHedgeFlag[0] = THOST_FTDC_HF_Speculation;

    order.VolumeTotalOriginal = 1;
    order.VolumeCondition = THOST_FTDC_VC_AV;
    order.MinVolume = 1;
    order.ForceCloseReason = THOST_FTDC_FCC_NotForceClose;
    order.IsAutoSuspend = 0;
    order.UserForceClose = 0;
    //cout<<"price: "<<price<<endl;
    if (price==0){
        order.OrderPriceType = THOST_FTDC_OPT_AnyPrice;
        order.LimitPrice = 0;
        order.TimeCondition = THOST_FTDC_TC_IOC;
    }else{
        order.OrderPriceType = THOST_FTDC_OPT_LimitPrice;
        order.LimitPrice = price;
        order.TimeCondition = THOST_FTDC_TC_GFD;
    }
    while (true)
    {
        int iResult = this->trApi->ReqOrderInsert(&order,this->getRequestID());
        if (!IsFlowControl(iResult))
        {
            break;
        }
        else
        {
            sleep(1);
        }
    }
}
void TraderSpi::reqUserLogin(){
    //if (NULL!=this->TradingDay)return;
    while (true)
    {
        int iResult = this->trApi->ReqUserLogin(&this->userReq,this->getRequestID());
        if (!IsFlowControl(iResult))
        {
            break;
        }
        else
        {
            sleep(1);
        }
    }
}



void TraderSpi::OnRtnTrade(CThostFtdcTradeField *pTrade) {
    cout<<"InstrumentID "<<pTrade->InstrumentID<<endl;
    cout<<"date "<< pTrade->TradeDate<<pTrade->TradeTime << endl;
    cout<<"Price "<< pTrade->Price << endl;
    cout<<"OrderRef "<< pTrade->OrderRef << endl;
    this->reqInvestorPosition(pTrade->InstrumentID);
}
void TraderSpi::OnRtnOrder(CThostFtdcOrderField *pOrder){

    //cout<<"order Ins "<< pOrder->InstrumentID << endl;
    //cout<<"order status "<< pOrder->OrderStatus << endl;
    //cout<<"order submit status "<< pOrder->OrderSubmitStatus << endl;
    //cout<<"order price "<< pOrder->LimitPrice << endl;

    //cout << "order VolumeTraded " << pOrder->VolumeTraded << endl;
    //cout << "order VolumeTotal "<< pOrder->VolumeTotal << endl;
}
