#include <iostream>
#include <unistd.h>
//#include <sys/stat.h>
//#include <string.h>
#include "traderspi.h"
#include <thread>
//#include <mutex>
//mutex mtx;
//map<string , CThostFtdcOrderField> mapOrder;

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
        const char * password,
        const char * passwordBak,
        const char * addr,
        const char * path):socketUnixServer(path){

    //this->Login = false;
    this->trApi = NULL;
    //this->TradingDay = NULL;
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

    this->mapstring["help"] = 999;
    this->mapstring["stop"] = 100;
    this->mapstring["Instrument"] = 1;
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
    this->mapstring["OrderInsert"] = 12;
    this->mapstring["OrderAction"] = 13;
}

void TraderSpi::routeHand(const char *data){

    //if (NULL==this->TradingDay){
    //    cout<<"trader:"<<data<<endl;
    //    return;
    //}
    //this->reqInstruments();
    char db[1024];
    strcpy(db,data);
    //cout<<"db "<<db<<endl;
    char *p;
    char sep[] = ",";
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
    case 1:{
        this->reqInstruments();
        break;
    }
    case 999:{
        this->help();
        break;
    }
    case 100:{
        this->stop();
        cout<<"stop"<<endl;
        break;
    }
    case 2:{
        this->setUserReg(str[1],str[2],str[3],str[4]);
        this->Addr = str[5];
        this->run();
        break;
    }
    case 3:{
        //this->setUserReg(str[1],str[2],str[3],str[4]);
        this->Addr = str[1];
        this->run();
        break;
    }
    case 4:{
        this->reqQrySettlementInfo();
        break;
    }
    case 5:{
        this->reqSettlementInfoConfirm();
        break;
    }
    case 6:{
        this->reqQrySettlementInfoConfirm();
        break;
    }
    case 7:{
        //cout<<"db:"<<db<<endl;
        this->reqTradingAccount();
        break;
    }
    case 8:{
        char * ins=NULL;
        if (i>1){
            ins = str[1];
        }
        this->reqInvestorPositionDetail(ins);
        break;
    }
    case 9:{
        char * ins=NULL;
        if (i>1){
            ins = str[1];
        }
        this->reqInvestorPosition(ins);
        break;
    }
    case 10:{
        if (i<5)break;
        double pr,pr_;
        //po = atof(str[4]);
        pr = atof(str[4]);
        pr_ = atof(str[5]);
        this->sendOrderOpen(str[1],str[2],str[3][0],pr,pr_);
        break;
    }
    case 11:{
        if (i<4)break;
        this->sendOrderClose(str[1],str[2],str[3][0],str[4][0]);
        break;
    }
    case 12:{
        if (i<6)break;
        this->sendOrderInsert(str[1],str[2],str[3],str[4][0],str[5][0],atof(str[6]),atof(str[7]));
        //this->sendOrderClose(str[1],str[2]);
        break;
    }
    case 13:{
        if (i<4)break;
        this->sendOrderAction(str[1],str[2],str[3],str[4]);
        break;
    }
    default:{
        printf("default %s %s end",data,str[0]);
        break;
    }
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
    cout<<"msg:InstrumentID 合约代码:"<<pInvestorPosition->InstrumentID<<endl;
    cout<<"msg:BrokerID 经纪公司代码:"<<pInvestorPosition->BrokerID<<endl;
    cout<<"msg:InvestorID 投资者代码:"<<pInvestorPosition->InvestorID<<endl;
    cout<<"msg:PosiDirection 持仓多空方向:"<<pInvestorPosition->PosiDirection<<endl;
    cout<<"msg:HedgeFlag 投机套保标志:"<<pInvestorPosition->HedgeFlag<<endl;
    cout<<"msg:PositionDate 持仓日期:"<<pInvestorPosition->PositionDate<<endl;
    cout<<"msg:YdPosition 上日持仓:"<<pInvestorPosition->YdPosition<<endl;
    cout<<"msg:Position 今日持仓:"<<pInvestorPosition->Position<<endl;
    cout<<"msg:LongFrozen 多头冻结:"<<pInvestorPosition->LongFrozen<<endl;
    cout<<"msg:ShortFrozen 空头冻结:"<<pInvestorPosition->ShortFrozen<<endl;
    cout<<"msg:LongFrozenAmount 开仓冻结金额:"<<pInvestorPosition->LongFrozenAmount<<endl;
    cout<<"msg:ShortFrozenAmount 开仓冻结金额:"<<pInvestorPosition->ShortFrozenAmount<<endl;
    cout<<"msg:OpenVolume 开仓量:"<<pInvestorPosition->OpenVolume<<endl;
    cout<<"msg:CloseVolume 平仓量:"<<pInvestorPosition->CloseVolume<<endl;
    cout<<"msg:OpenAmount 开仓金额:"<<pInvestorPosition->OpenAmount<<endl;
    cout<<"msg:CloseAmount 平仓金额:"<<pInvestorPosition->CloseAmount<<endl;
    cout<<"msg:PositionCost 持仓成本:"<<pInvestorPosition->PositionCost<<endl;
    cout<<"msg:PreMargin 上次占用的保证金:"<<pInvestorPosition->PreMargin<<endl;
    cout<<"msg:UseMargin 占用的保证金:"<<pInvestorPosition->UseMargin<<endl;
    cout<<"msg:FrozenMargin 冻结的保证金:"<<pInvestorPosition->FrozenMargin<<endl;
    cout<<"msg:FrozenCash 冻结的资金:"<<pInvestorPosition->FrozenCash<<endl;
    cout<<"msg:FrozenCommission 冻结的手续费:"<<pInvestorPosition->FrozenCommission<<endl;
    cout<<"msg:CashIn 资金差额:"<<pInvestorPosition->CashIn<<endl;
    cout<<"msg:Commission 手续费:"<<pInvestorPosition->Commission<<endl;
    cout<<"msg:CloseProfit 平仓盈亏:"<<pInvestorPosition->CloseProfit<<endl;
    cout<<"msg:PositionProfit 持仓盈亏:"<<pInvestorPosition->PositionProfit<<endl;
    cout<<"msg:PreSettlementPrice 上次结算价:"<<pInvestorPosition->PreSettlementPrice<<endl;
    cout<<"msg:SettlementPrice 本次结算价:"<<pInvestorPosition->SettlementPrice<<endl;
    cout<<"msg:TradingDay 交易日:"<<pInvestorPosition->TradingDay<<endl;
    cout<<"msg:SettlementID 结算编号:"<<pInvestorPosition->SettlementID<<endl;
    cout<<"msg:OpenCost 开仓成本:"<<pInvestorPosition->OpenCost<<endl;
    cout<<"msg:ExchangeMargin 交易所保证金:"<<pInvestorPosition->ExchangeMargin<<endl;
    cout<<"msg:CombPosition 组合成交形成的持仓:"<<pInvestorPosition->CombPosition<<endl;
    cout<<"msg:CombLongFrozen 组合多头冻结:"<<pInvestorPosition->CombLongFrozen<<endl;
    cout<<"msg:CombShortFrozen 组合空头冻结:"<<pInvestorPosition->CombShortFrozen<<endl;
    cout<<"msg:CloseProfitByDate 逐日盯市平仓盈亏:"<<pInvestorPosition->CloseProfitByDate<<endl;
    cout<<"msg:CloseProfitByTrade 逐笔对冲平仓盈亏:"<<pInvestorPosition->CloseProfitByTrade<<endl;
    cout<<"msg:TodayPosition 今日持仓:"<<pInvestorPosition->TodayPosition<<endl;
    cout<<"msg:MarginRateByMoney 保证金率:"<<pInvestorPosition->MarginRateByMoney<<endl;
    cout<<"msg:MarginRateByVolume 保证金率(按手数):"<<pInvestorPosition->MarginRateByVolume<<endl;
    cout<<"msg:StrikeFrozen 执行冻结:"<<pInvestorPosition->StrikeFrozen<<endl;
    cout<<"msg:StrikeFrozenAmount 执行冻结金额:"<<pInvestorPosition->StrikeFrozenAmount<<endl;
    cout<<"msg:AbandonFrozen 放弃执行冻结:"<<pInvestorPosition->AbandonFrozen<<endl;
    cout<<"msg:ExchangeID 交易所代码:"<<pInvestorPosition->ExchangeID<<endl;
    cout<<"msg:YdStrikeFrozen 执行冻结的昨仓:"<<pInvestorPosition->YdStrikeFrozen<<endl;
    cout<<"msg:InvestUnitID 投资单元代码:"<<pInvestorPosition->InvestUnitID<<endl;
    //cout<<"msg:PositionCostOffset 大商所持仓成本差值，只有大商所使用:"<<pInvestorPosition->PositionCostOffset<endl;


}

void TraderSpi::OnRspQryTradingAccount(
        CThostFtdcTradingAccountField *pTradingAccount,
        CThostFtdcRspInfoField *pRspInfo,
        int nRequestID,
        bool bIsLast) {
    //return;
    //cout<<"onRspTA"<<endl;
    if (pRspInfo && pRspInfo->ErrorID!=0){
        cout<<pRspInfo->ErrorMsg<<endl;
        return;
    }

    if (!bIsLast)return;

    if (!pTradingAccount)return;
    char db[1024];
    sprintf(db,"ta %lf",pTradingAccount->Deposit);
    cout<<db<<endl;
    this->send(db);

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
    cout << "InvestorPositionDetail "\
         << pInvestorPositionDetail->BrokerID <<","\
         << pInvestorPositionDetail->OpenDate <<","\
         << pInvestorPositionDetail->TradingDay << ","\
         << pInvestorPositionDetail->InstrumentID << ","\
         << pInvestorPositionDetail->ExchangeID << ","\
         << pInvestorPositionDetail->Direction << endl;

}

void TraderSpi::OnErrRtnOrderInsert(
		CThostFtdcInputOrderField *pInputOrder, 
		CThostFtdcRspInfoField *pRspInfo) {

    if (pRspInfo && 0!=pRspInfo->ErrorID){

        cout<<"orderCancel "<<pInputOrder->InstrumentID<<" "<<pInputOrder->OrderRef<<endl;
        //cout<<"err order insert "<<pRspInfo->ErrorMsg<<pInputOrder->CombOffsetFlag[0]<<endl;
        return;
    }
    //cout<<"order err"<<pInputOrder->OrderRef<<endl;

}
void TraderSpi::OnRspOrderInsert(
        CThostFtdcInputOrderField *pInputOrder,
        CThostFtdcRspInfoField *pRspInfo,
        int nRequestID,
        bool bIsLast){
    if (pRspInfo && (0 != pRspInfo->ErrorID)){

	    
        //cout<<"orderCancel "<<pInputOrder.->InstrumentID<<" "<<pInputOrder.->OrderRef<<endl;
        cout<<"err order "<<pRspInfo->ErrorMsg<<pInputOrder->CombOffsetFlag[0]<<endl;
        return;
    }
    //cout<<"order insert "<<pInputOrder->OrderRef<<endl;
}
void TraderSpi::help(){
    map<string , int>::iterator iter;
    for(iter = mapstring.begin(); iter != mapstring.end(); iter++)
          //cout<<iter->first<<endl;
          cout<<"help:"<<iter->first<<endl;

}

void TraderSpi::stop(){
    if (this->trApi == NULL) return;
    cout<<"stop ok 3"<<endl;
    this->trApi->RegisterSpi(NULL);
    cout<<"stop ok 2"<<endl;
    this->trApi->Release();
    cout<<"stop ok 1"<<endl;
    this->trApi = NULL;
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
    //this->mapInstrument[pInstrument->InstrumentID] = *pInstrument;
    char db[8196] = "ins ";
    sprintf(db,
            "ins InstrumentID:%s,ExchangeID:%s,InstrumentName:%s,PriceTick:%lf,CreateDate:%s,OpenDate:%s,ExpireDate:%s,StartDelivDate:%s,EndDelivDate:%s,IsTrading:%d}",
            pInstrument->InstrumentID,
            pInstrument->ExchangeID,
            pInstrument->InstrumentName,
            pInstrument->PriceTick,
            pInstrument->CreateDate,
            pInstrument->OpenDate,
            pInstrument->ExpireDate,
            pInstrument->StartDelivDate,
            pInstrument->EndDelivDate,
            pInstrument->IsTrading);
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
        //this->Login = true;
        this->frontID = pRspUserLogin->FrontID;
        this->sessionID = pRspUserLogin->SessionID;

        char trading[20]="TDay ";
        strcat(trading,this->trApi->GetTradingDay());
        this->send(trading);
        //strcpy(this->TradingDay,this->trApi->GetTradingDay());
        cout <<"Td connected "<< trading << endl;
        cout <<"frontID:"<< this->frontID << endl;
        cout <<"sessionID:"<< this->sessionID << endl;
        //this->reqInstruments();
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

void TraderSpi::run(){
    if (this->trApi != NULL) return;
    this->trApi = CThostFtdcTraderApi::CreateFtdcTraderApi(this->path);
    this->trApi->RegisterSpi(this);
    char _addr[1024];
    strcpy(_addr,Addr);
    //cout<<_addr<<endl;
    this->trApi->SubscribePublicTopic(THOST_TERT_QUICK);				// 注册公有流
    this->trApi->SubscribePrivateTopic(THOST_TERT_QUICK);
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

    //cout<<"OnRspQrySettlementInfo:"<<bIsLast<<endl;
    if (!bIsLast)return;
    //char msg[8192];


    if (!pSettlementInfo) return;
    //cout<<pSettlementInfo->TradingDay<<endl;
    //cout<<pSettlementInfo->SequenceNo<< endl;
    //cout<<pSettlementInfo->SettlementID<< endl;
    //cout<<pSettlementInfo->Content<< endl;
    char msg[1024];
    sprintf(msg,
            "msg TradingDay|%s SequenceNo|%d SettlementID|%d Content|%s",
            pSettlementInfo->TradingDay,
            pSettlementInfo->SequenceNo,
            pSettlementInfo->SettlementID,
            pSettlementInfo->Content);
    cout<<msg<<endl;
    this->send(msg);
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

    //cout<<"OnRspQrySettlementInfoConfirm:"<<bIsLast<<endl;
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
    //if (!this->Login)return;
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
    //if (!this->Login)return;
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
    //if (!this->Login)return;
    CThostFtdcQrySettlementInfoField pQrySettlementInfo;
    memset(&pQrySettlementInfo,0,sizeof(pQrySettlementInfo));
    strcpy(pQrySettlementInfo.BrokerID,this->userReq.BrokerID);
    strcpy(pQrySettlementInfo.AccountID,this->userReq.UserID);
    strcpy(pQrySettlementInfo.InvestorID,this->userReq.UserID);
    while (true)
    {
        int iResult = this->trApi->ReqQrySettlementInfo(&pQrySettlementInfo,this->getRequestID());
        cout<<iResult<<endl;
        char msg[1024];
        sprintf(msg,"req:%d",iResult);
        this->send(msg);
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

    //if (!this->Login)return;
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

    //cout<<"login:"<<this->Login<<endl;
    //if (!this->Login)return;
    //cout<<"req trading Account"<<endl;
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
    //if (!this->Login)return;
    CThostFtdcQryInvestorPositionField pQryInvestorPosition;
    memset(&pQryInvestorPosition,0,sizeof(pQryInvestorPosition));
    strcpy(pQryInvestorPosition.InvestorID,this->userReq.UserID);
    strcpy(pQryInvestorPosition.BrokerID,this->userReq.BrokerID);
    if (ins!=NULL){
        cout<<"info "<<ins<<endl;
        strcpy(pQryInvestorPosition.InstrumentID,ins);
    }
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

    //if (!this->Login)return;
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

void TraderSpi::sendOrderClose(const char * ins, const char * ExchangeID, const char dis, const char type){

    //if (!this->Login)return;
    //CThostFtdcInstrumentField insinfo = this->mapInstrument[ins];
    cout<<"close "<<ins<<endl;
    CThostFtdcInputOrderField order;
    memset(&order,0,sizeof(order));
    strcpy(order.BrokerID,this->userReq.BrokerID);
    strcpy(order.InvestorID,this->userReq.UserID);
    strcpy(order.InstrumentID,ins);
    //strcpy(order.OrderRef,OrderRef);
    strcpy(order.UserID,this->userReq.UserID);
    strcpy(order.ExchangeID,ExchangeID);
    order.ContingentCondition =THOST_FTDC_CC_Immediately;
    order.Direction = dis;

    //order.CombOffsetFlag[0] = THOST_FTDC_OF_CloseToday;
    order.CombOffsetFlag[0] = type;
    //order.CombOffsetFlag[3] = THOST_FTDC_OF_CloseYesterday;
    order.CombHedgeFlag[0] = THOST_FTDC_HF_Speculation;

    order.VolumeTotalOriginal = 1;
    order.VolumeCondition = THOST_FTDC_VC_AV;
    order.MinVolume = 1;
    order.ForceCloseReason = THOST_FTDC_FCC_NotForceClose;
    order.IsAutoSuspend = 0;
    order.UserForceClose = 0;

    order.OrderPriceType = THOST_FTDC_OPT_AnyPrice;
    order.LimitPrice = 0;
    order.TimeCondition = THOST_FTDC_TC_IOC;

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

void TraderSpi::sendOrderOpen(
        const char *ins,
        const char *ExchangeID,
        const char dir,
        const double price,
        const double stopPrice){
    cout<<"open "<<ins<<endl;
    CThostFtdcInputOrderField order;
    memset(&order,0,sizeof(order));
    strcpy(order.BrokerID,this->userReq.BrokerID);
    strcpy(order.InvestorID,this->userReq.UserID);
    strcpy(order.InstrumentID,ins);
    //strcpy(order.OrderRef,orderRef);
    strcpy(order.UserID,this->userReq.UserID);
    strcpy(order.ExchangeID,ExchangeID);
    order.ContingentCondition =THOST_FTDC_CC_Immediately;
    //order.ContingentCondition = THOST_FTDC_CC_Touch;
    order.StopPrice=stopPrice;
    order.Direction = dir;
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

void TraderSpi::sendOrderInsert(
        const char *ins,
        const char *ExchangeID,
        const char *OrderRef,
        const char setFlag,
        const char dis,
        const double price,
        const double stopPrice){

    CThostFtdcInputOrderField order;
    memset(&order,0,sizeof(order));
    strcpy(order.BrokerID,this->userReq.BrokerID);
    strcpy(order.InvestorID,this->userReq.UserID);
    strcpy(order.InstrumentID,ins);
    strcpy(order.UserID,this->userReq.UserID);
    strcpy(order.ExchangeID,ExchangeID);
    //strcpy(order.OrderRef,OrderRef);
    order.ContingentCondition =THOST_FTDC_CC_Immediately;
    order.Direction = dis;
    order.CombOffsetFlag[0] = setFlag;
    order.CombHedgeFlag[0] = THOST_FTDC_HF_Speculation;

    order.VolumeTotalOriginal = 1;
    order.VolumeCondition = THOST_FTDC_VC_AV;
    order.MinVolume = 1;
    order.ForceCloseReason = THOST_FTDC_FCC_NotForceClose;
    order.IsAutoSuspend = 0;
    order.UserForceClose = 0;
    //cout<<"price: "<<price<<endl;
    //if (price==0){
    //    order.OrderPriceType = THOST_FTDC_OPT_AnyPrice;
    //    order.LimitPrice = 0;
    //    order.TimeCondition = THOST_FTDC_TC_IOC;
    //}else{
    order.OrderPriceType = THOST_FTDC_OPT_LimitPrice;
    order.LimitPrice = price;
    order.TimeCondition = THOST_FTDC_TC_GFD;
    if (stopPrice>0)order.StopPrice = stopPrice;
    //order.TimeCondition = THOST_FTDC_TC_IOC;
    //}
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

void TraderSpi::sendOrderAction(
        const char *ins,
        const char *ExchangeID,
        const char *OrderRef,
        const char *OrderSysID
        ){

    CThostFtdcInputOrderActionField action;
    memset(&action,0,sizeof(action));
    strcpy(action.BrokerID,this->userReq.BrokerID);
    strcpy(action.InvestorID,this->userReq.UserID);
    strcpy(action.InstrumentID,ins);
    strcpy(action.UserID,this->userReq.UserID);
    strcpy(action.ExchangeID,ExchangeID);
    action.ActionFlag  = THOST_FTDC_AF_Delete;
    action.FrontID = this->frontID;
    action.SessionID = this->sessionID;
    //TThostFtdcOrderRefType Oref = this->mapOrder[ins]
    //strcpy(action.OrderRef,this->mapOrder[ins].OrderRef);
    strcpy(action.OrderRef,OrderRef);
    strcpy(action.OrderSysID,OrderSysID);
    //string _ins = action.InstrumentID;
    //mtx.lock();
    //if(mapOrder.find(_ins) == mapOrder.end()){
    //	cout<<"findNot "<< _ins<<endl;
    //    map<string , CThostFtdcOrderField>::iterator iter;  
    //    for(iter = mapOrder.begin(); iter != mapOrder.end(); iter++)  
    //        cout<<"check "<<iter->first<<' '<<iter->second.OrderRef<<endl;  
    //}else{
    //	cout<<"orderRef "<< mapOrder[_ins].OrderRef << endl;
    //	strcpy(action.OrderRef,mapOrder[_ins].OrderRef);
    //}
    //mtx.unlock();
    //cout << action.OrderRef << "-" << OrderRef << "-" << ExchangeID <<endl;

    while (true)
    {
        int iResult = this->trApi->ReqOrderAction(&action,this->getRequestID());
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
    cout<<"trade "\
       << pTrade->InstrumentID << " " \
       << pTrade->Price << " " \
       << pTrade->OffsetFlag << " " \
       << pTrade->TradeDate << "T" \
       << pTrade->TradeTime << " " \
       << pTrade->OrderRef << " " \
       << endl;
    //this->reqInvestorPosition(pTrade->InstrumentID);
}
void TraderSpi::OnRtnOrder(CThostFtdcOrderField *pOrder){

    if (pOrder->CombOffsetFlag[0] == THOST_FTDC_OF_Open){
   	 //if (pOrder->OrderStatus == '3' || pOrder->OrderStatus=='4'){
   	 //    this->sendOrderAction(pOrder->InstrumentID,pOrder->ExchangeID,pOrder->OrderRef);
	 //    return;
   	 //}
   	 //    cout<<"orderWait "<<pOrder->InstrumentID<<" "<<pOrder->OrderRef<<endl;
   	 //}else if (pOrder->OrderStatus=='5'){
   	 if (pOrder->OrderStatus=='5'){
   	     cout<<"orderCancel "<<pOrder->InstrumentID<<"-"<<pOrder->OrderRef<<endl;
   	 //}else if (pOrder->OrderStatus!='0'){
   	 }else{
   	     //cout<<"orderWait "<<pOrder->InstrumentID<<"-"<<pOrder->OrderSysID<<endl;
   	     cout <<"orderWait "<<pOrder->InstrumentID << "-";
	     cout <<pOrder->OrderRef << "-" << pOrder->OrderSysID << endl;
	     //iter = this->mapOrder.find(pOrder->InstrumentID);
	     //string ins = pOrder->InstrumentID;
	     //mtx.lock();
    	     //if(mapOrder.find(ins)==mapOrder.end()){
	     //  mapOrder.insert(map<string,CThostFtdcOrderField>::value_type(ins,*pOrder));
	     //  cout<<"insert "<< ins << endl;
	     //}else{
	     //  mapOrder[ins]=*pOrder;
	     //  cout<<"update "<< ins << endl;
	     //}
	     //mtx.unlock();
   	 }
    }
    //cout<<"msg:InstrumentID 合约代码:"<< pOrder->InstrumentID << endl;
    //cout<<"msg:order OrderRef "<< pOrder->OrderRef << endl;
    //cout<<"msg:order status "<< pOrder->OrderStatus << endl;
    //cout<<"msg:order submit status "<< pOrder->OrderSubmitStatus << endl;
    //cout<<"msg:order price "<< pOrder->UpdateTime << endl;
    //cout<<"msg:order price "<< pOrder->LimitPrice << endl;
    //cout<<"StatusMsg:"<< pOrder->StatusMsg << endl;
    //cout<<"Order--------------------------" << endl;
}

void TraderSpi::OnRspOrderAction(
        CThostFtdcInputOrderActionField *pInputOrderAction,
        CThostFtdcRspInfoField *pRspInfo, int nRequestID, bool bIsLast){

    if (pRspInfo && pRspInfo->ErrorID!=0){
        cout<<"action info ";
        cout<< pRspInfo->ErrorMsg<<endl;
        cout<< pInputOrderAction->InstrumentID <<",";
        cout<< pInputOrderAction->ExchangeID <<",";
        cout<< this->frontID <<",";
        cout<< this->sessionID <<",";
	cout<< pInputOrderAction->OrderRef <<endl;
        return;
    }
    //if (!bIsLast)return;

}
