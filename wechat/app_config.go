package wechat

type AppConfig struct {
	SubPackages []struct {
		Root    string `json:"root"`
		Plugins struct {
			LivePlayerPlugin struct {
				Version    string `json:"version"`
				Provider   string `json:"provider"`
				Subpackage string `json:"subpackage"`
			} `json:"live-player-plugin"`
		} `json:"plugins,omitempty"`
	} `json:"subPackages"`
	EntryPagePath string   `json:"entryPagePath"`
	Pages         []string `json:"pages"`
	Page          struct {
		PagesIndexIndexHtml struct {
			Window struct {
				NavigationBarBackgroundColor string `json:"navigationBarBackgroundColor"`
				NavigationStyle              string `json:"navigationStyle"`
				OnReachBottomDistance        int    `json:"onReachBottomDistance"`
			} `json:"window"`
		} `json:"pages/index/index.html"`
		PagesEstateTopicIndexHtml struct {
			Window struct {
				NavigationBarBackgroundColor string `json:"navigationBarBackgroundColor"`
				NavigationStyle              string `json:"navigationStyle"`
				OnReachBottomDistance        int    `json:"onReachBottomDistance"`
			} `json:"window"`
		} `json:"pages/estateTopic/index.html"`
		PagesArticleIndexHtml struct {
			Window struct {
			} `json:"window"`
		} `json:"pages/article/index.html"`
		PagesBuildPageIndexHtml struct {
			Window struct {
				NavigationStyle string `json:"navigationStyle"`
			} `json:"window"`
		} `json:"pages/buildPage/index.html"`
		PagesDetailLabelIndexHtml struct {
			Window struct {
				NavigationStyle string `json:"navigationStyle"`
			} `json:"window"`
		} `json:"pages/detailLabel/index.html"`
		PagesLogsIndexHtml struct {
			Window struct {
				NavigationStyle string `json:"navigationStyle"`
			} `json:"window"`
		} `json:"pages/logs/index.html"`
		PagesMessageIndexHtml struct {
			Window struct {
				NavigationStyle string `json:"navigationStyle"`
			} `json:"window"`
		} `json:"pages/message/index.html"`
		PagesMyChatsIndexHtml struct {
			Window struct {
				NavigationStyle       string `json:"navigationStyle"`
				EnablePullDownRefresh bool   `json:"enablePullDownRefresh"`
				DisableScroll         bool   `json:"disableScroll"`
			} `json:"window"`
		} `json:"pages/myChats/index.html"`
		PagesQwaIndexHtml struct {
			Window struct {
				NavigationStyle string `json:"navigationStyle"`
			} `json:"window"`
		} `json:"pages/qwa/index.html"`
		PagesShHouseIndexHtml struct {
			Window struct {
				NavigationStyle string `json:"navigationStyle"`
			} `json:"window"`
		} `json:"pages/shHouse/index.html"`
		PackTab1PKIndexHtml struct {
			Window struct {
				NavigationStyle string `json:"navigationStyle"`
			} `json:"window"`
		} `json:"packTab1/PK/index.html"`
		PackTab1AuthorizationPhoneIndexHtml struct {
			Window struct {
			} `json:"window"`
		} `json:"packTab1/authorizationPhone/index.html"`
		PackTab1SchoolMapIndexHtml struct {
			Window struct {
				NavigationStyle       string `json:"navigationStyle"`
				EnablePullDownRefresh bool   `json:"enablePullDownRefresh"`
				DisableScroll         bool   `json:"disableScroll"`
			} `json:"window"`
		} `json:"packTab1/schoolMap/index.html"`
		PackTab1MapChooseIndexHtml struct {
			Window struct {
				NavigationStyle string `json:"navigationStyle"`
				DisableScroll   bool   `json:"disableScroll"`
			} `json:"window"`
		} `json:"packTab1/mapChoose/index.html"`
		PackTab1MapAreaAllIndexHtml struct {
			Window struct {
				NavigationStyle string `json:"navigationStyle"`
				DisableScroll   bool   `json:"disableScroll"`
			} `json:"window"`
		} `json:"packTab1/mapAreaAll/index.html"`
		PackTab1OrderPageIndexHtml struct {
			Window struct {
				NavigationStyle string `json:"navigationStyle"`
			} `json:"window"`
		} `json:"packTab1/orderPage/index.html"`
		PackTab1ApartmentIndexHtml struct {
			Window struct {
				NavigationStyle string `json:"navigationStyle"`
			} `json:"window"`
		} `json:"packTab1/Apartment/index.html"`
		PackTab1AreaSortIndexHtml struct {
			Window struct {
				NavigationStyle string `json:"navigationStyle"`
			} `json:"window"`
		} `json:"packTab1/areaSort/index.html"`
		PackTab1MapAreaIndexHtml struct {
			Window struct {
				NavigationStyle       string `json:"navigationStyle"`
				EnablePullDownRefresh bool   `json:"enablePullDownRefresh"`
			} `json:"window"`
		} `json:"packTab1/mapArea/index.html"`
		PackTab1MapLocateIndexHtml struct {
			Window struct {
				NavigationStyle string `json:"navigationStyle"`
			} `json:"window"`
		} `json:"packTab1/mapLocate/index.html"`
		PackTab1MoreLocaleIndexHtml struct {
			Window struct {
				NavigationStyle       string `json:"navigationStyle"`
				EnablePullDownRefresh bool   `json:"enablePullDownRefresh"`
			} `json:"window"`
		} `json:"packTab1/moreLocale/index.html"`
		PackTab1PKreaultIndexHtml struct {
			Window struct {
				NavigationStyle string `json:"navigationStyle"`
			} `json:"window"`
		} `json:"packTab1/PKreault/index.html"`
		PackTab1TalkAskIndexHtml struct {
			Window struct {
				NavigationStyle string `json:"navigationStyle"`
			} `json:"window"`
		} `json:"packTab1/talkAsk/index.html"`
		PackTab1TalkAskAndQuestIndexHtml struct {
			Window struct {
				NavigationStyle string `json:"navigationStyle"`
			} `json:"window"`
		} `json:"packTab1/talkAskAndQuest/index.html"`
		PackTab1TalkMoreIndexHtml struct {
			Window struct {
				NavigationStyle string `json:"navigationStyle"`
			} `json:"window"`
		} `json:"packTab1/talkMore/index.html"`
		PackTab1TalkTopicIndexHtml struct {
			Window struct {
			} `json:"window"`
		} `json:"packTab1/talkTopic/index.html"`
		PackTab1TalkUserIndexHtml struct {
			Window struct {
				NavigationStyle       string `json:"navigationStyle"`
				EnablePullDownRefresh bool   `json:"enablePullDownRefresh"`
			} `json:"window"`
		} `json:"packTab1/talkUser/index.html"`
		PackTab1CallHistoryIndexHtml struct {
			Window struct {
				NavigationStyle string `json:"navigationStyle"`
			} `json:"window"`
		} `json:"packTab1/callHistory/index.html"`
		PackTab1BaseAllInfoIndexHtml struct {
			Window struct {
				NavigationStyle string `json:"navigationStyle"`
			} `json:"window"`
		} `json:"packTab1/baseAllInfo/index.html"`
		PackTab1DengjiSearchIndexHtml struct {
			Window struct {
				NavigationStyle string `json:"navigationStyle"`
			} `json:"window"`
		} `json:"packTab1/dengjiSearch/index.html"`
		PackTab1DesDetailIndexHtml struct {
			Window struct {
				NavigationStyle string `json:"navigationStyle"`
			} `json:"window"`
		} `json:"packTab1/desDetail/index.html"`
		PackTab1OneLandInfoIndexHtml struct {
			Window struct {
				NavigationStyle string `json:"navigationStyle"`
			} `json:"window"`
		} `json:"packTab1/oneLandInfo/index.html"`
		PackTab1ReplyMesIndexHtml struct {
			Window struct {
			} `json:"window"`
		} `json:"packTab1/replyMes/index.html"`
		PackTab1TouchIndexHtml struct {
			Window struct {
				NavigationStyle string `json:"navigationStyle"`
			} `json:"window"`
		} `json:"packTab1/touch/index.html"`
		PackTab1RealAllImagesIndexHtml struct {
			Window struct {
				NavigationStyle       string `json:"navigationStyle"`
				EnablePullDownRefresh bool   `json:"enablePullDownRefresh"`
			} `json:"window"`
		} `json:"packTab1/realAllImages/index.html"`
		PackTab1ChangeUserIndexHtml struct {
			Window struct {
				NavigationStyle string `json:"navigationStyle"`
			} `json:"window"`
		} `json:"packTab1/changeUser/index.html"`
		PackTab1TestPageIndexHtml struct {
			Window struct {
				NavigationStyle string `json:"navigationStyle"`
			} `json:"window"`
		} `json:"packTab1/testPage/index.html"`
		PackTab1TestPage2IndexHtml struct {
			Window struct {
				NavigationStyle       string `json:"navigationStyle"`
				EnablePullDownRefresh bool   `json:"enablePullDownRefresh"`
				DisableScroll         bool   `json:"disableScroll"`
			} `json:"window"`
		} `json:"packTab1/testPage2/index.html"`
		PackTab1TestPage3IndexHtml struct {
			Window struct {
				NavigationStyle       string `json:"navigationStyle"`
				EnablePullDownRefresh bool   `json:"enablePullDownRefresh"`
				DisableScroll         bool   `json:"disableScroll"`
			} `json:"window"`
		} `json:"packTab1/testPage3/index.html"`
		PackTab1TestResultIndexHtml struct {
			Window struct {
				NavigationStyle string `json:"navigationStyle"`
			} `json:"window"`
		} `json:"packTab1/testResult/index.html"`
		PackTab1LoginCodeIndexHtml struct {
			Window struct {
				NavigationStyle string `json:"navigationStyle"`
			} `json:"window"`
		} `json:"packTab1/loginCode/index.html"`
		PackTab1IdentityIndexHtml struct {
			Window struct {
				NavigationStyle string `json:"navigationStyle"`
			} `json:"window"`
		} `json:"packTab1/identity/index.html"`
		PackTab1DetailInfoIndexHtml struct {
			Window struct {
				NavigationStyle string `json:"navigationStyle"`
			} `json:"window"`
		} `json:"packTab1/detailInfo/index.html"`
		PackTab1NeedRealIndexHtml struct {
			Window struct {
				NavigationStyle string `json:"navigationStyle"`
			} `json:"window"`
		} `json:"packTab1/needReal/index.html"`
		PackTab1SecHouseIndexHtml struct {
			Window struct {
			} `json:"window"`
		} `json:"packTab1/secHouse/index.html"`
		PackTab1AssembleIndexHtml struct {
			Window struct {
				NavigationStyle string `json:"navigationStyle"`
			} `json:"window"`
		} `json:"packTab1/assemble/index.html"`
		PackTab1ToparticularsIndexHtml struct {
			Window struct {
				NavigationStyle string `json:"navigationStyle"`
			} `json:"window"`
		} `json:"packTab1/toparticulars/index.html"`
		PackTab1TransactionIndexHtml struct {
			Window struct {
			} `json:"window"`
		} `json:"packTab1/transaction/index.html"`
		PackTab1InvestIndexHtml struct {
			Window struct {
			} `json:"window"`
		} `json:"packTab1/invest/index.html"`
		PackTab1CompanyMapIndexHtml struct {
			Window struct {
				NavigationBarBackgroundColor string `json:"navigationBarBackgroundColor"`
			} `json:"window"`
		} `json:"packTab1/companyMap/index.html"`
		PackTab1TopicMoreIndexHtml struct {
			Window struct {
			} `json:"window"`
		} `json:"packTab1/topicMore/index.html"`
		PackTab1ThirtyScoreIndexHtml struct {
			Window struct {
				NavigationStyle string `json:"navigationStyle"`
			} `json:"window"`
		} `json:"packTab1/thirtyScore/index.html"`
		PackTab1UserProfileIndexHtml struct {
			Window struct {
				NavigationStyle string `json:"navigationStyle"`
			} `json:"window"`
		} `json:"packTab1/userProfile/index.html"`
		PackTab1CollectIndexHtml struct {
			Window struct {
				NavigationStyle string `json:"navigationStyle"`
			} `json:"window"`
		} `json:"packTab1/collect/index.html"`
		PackTab1ErHouseIndexHtml struct {
			Window struct {
			} `json:"window"`
		} `json:"packTab1/erHouse/index.html"`
		PackTab1MoreStocksIndexHtml struct {
			Window struct {
			} `json:"window"`
		} `json:"packTab1/moreStocks/index.html"`
		PackTab1WithFriendsIndexHtml struct {
			Window struct {
			} `json:"window"`
		} `json:"packTab1/withFriends/index.html"`
		PackTab1EstateMapIndexHtml struct {
			Window struct {
				NavigationStyle       string `json:"navigationStyle"`
				EnablePullDownRefresh bool   `json:"enablePullDownRefresh"`
				DisableScroll         bool   `json:"disableScroll"`
			} `json:"window"`
		} `json:"packTab1/EstateMap/index.html"`
		PackTab1PlateMapIndexHtml struct {
			Window struct {
				NavigationStyle string `json:"navigationStyle"`
			} `json:"window"`
		} `json:"packTab1/plateMap/index.html"`
		PackTab1RegisterPageIndexHtml struct {
			Window struct {
				NavigationStyle string `json:"navigationStyle"`
			} `json:"window"`
		} `json:"packTab1/RegisterPage/index.html"`
		PackTab1SubwayMapPageIndexHtml struct {
			Window struct {
				NavigationStyle       string `json:"navigationStyle"`
				EnablePullDownRefresh bool   `json:"enablePullDownRefresh"`
				DisableScroll         bool   `json:"disableScroll"`
			} `json:"window"`
		} `json:"packTab1/SubwayMapPage/index.html"`
		PackTab2MortgageRateIndexHtml struct {
			Window struct {
				NavigationStyle string `json:"navigationStyle"`
			} `json:"window"`
		} `json:"packTab2/MortgageRate/index.html"`
		PackTab2AddLiveIndexHtml struct {
			Window struct {
				NavigationStyle string `json:"navigationStyle"`
			} `json:"window"`
		} `json:"packTab2/addLive/index.html"`
		PackTab2QwaTitleIndexHtml struct {
			Window struct {
				NavigationStyle string `json:"navigationStyle"`
			} `json:"window"`
		} `json:"packTab2/qwaTitle/index.html"`
		PackTab2IntelligenceIndexHtml struct {
			Window struct {
				NavigationStyle string `json:"navigationStyle"`
			} `json:"window"`
		} `json:"packTab2/intelligence/index.html"`
		PackTab2SelfTestIndexHtml struct {
			Window struct {
				NavigationStyle       string `json:"navigationStyle"`
				EnablePullDownRefresh bool   `json:"enablePullDownRefresh"`
				DisableScroll         bool   `json:"disableScroll"`
			} `json:"window"`
		} `json:"packTab2/selfTest/index.html"`
		PackTab2SearchResultIndexHtml struct {
			Window struct {
				NavigationStyle string `json:"navigationStyle"`
			} `json:"window"`
		} `json:"packTab2/searchResult/index.html"`
		PackTab2TalentsIndexHtml struct {
			Window struct {
				NavigationStyle       string `json:"navigationStyle"`
				EnablePullDownRefresh bool   `json:"enablePullDownRefresh"`
				DisableScroll         bool   `json:"disableScroll"`
			} `json:"window"`
		} `json:"packTab2/talents/index.html"`
		PackTab2IntegralIndexHtml struct {
			Window struct {
				NavigationStyle       string `json:"navigationStyle"`
				EnablePullDownRefresh bool   `json:"enablePullDownRefresh"`
				DisableScroll         bool   `json:"disableScroll"`
			} `json:"window"`
		} `json:"packTab2/integral/index.html"`
		PackTab2HotRealIndexHtml struct {
			Window struct {
				NavigationStyle string `json:"navigationStyle"`
			} `json:"window"`
		} `json:"packTab2/hotReal/index.html"`
		PackTab2FindBarIndexHtml struct {
			Window struct {
				NavigationStyle string `json:"navigationStyle"`
			} `json:"window"`
		} `json:"packTab2/findBar/index.html"`
		PackTab2MyClientIndexHtml struct {
			Window struct {
				NavigationStyle string `json:"navigationStyle"`
			} `json:"window"`
		} `json:"packTab2/myClient/index.html"`
		PackTab2RealDetailIndexHtml struct {
			Window struct {
				NavigationStyle string `json:"navigationStyle"`
			} `json:"window"`
		} `json:"packTab2/realDetail/index.html"`
		PackTab2RealNewsIndexHtml struct {
			Window struct {
			} `json:"window"`
		} `json:"packTab2/realNews/index.html"`
		PackTab2EditIndexHtml struct {
			Window struct {
				NavigationStyle string `json:"navigationStyle"`
			} `json:"window"`
		} `json:"packTab2/edit/index.html"`
		PackTab2EditListIndexHtml struct {
			Window struct {
			} `json:"window"`
		} `json:"packTab2/editList/index.html"`
		PackTab2CardScoreIndexHtml struct {
			Window struct {
				NavigationStyle string `json:"navigationStyle"`
			} `json:"window"`
		} `json:"packTab2/cardScore/index.html"`
		PackTab2TalkShareIndexHtml struct {
			Window struct {
				NavigationStyle string `json:"navigationStyle"`
			} `json:"window"`
		} `json:"packTab2/talkShare/index.html"`
		PackTab2GetMoreTurnsIndexHtml struct {
			Window struct {
				NavigationStyle string `json:"navigationStyle"`
			} `json:"window"`
		} `json:"packTab2/getMoreTurns/index.html"`
		PackTab2GetTurnsIndexHtml struct {
			Window struct {
				NavigationStyle string `json:"navigationStyle"`
			} `json:"window"`
		} `json:"packTab2/getTurns/index.html"`
		PackTab2AddScoreIndexHtml struct {
			Window struct {
				NavigationStyle string `json:"navigationStyle"`
			} `json:"window"`
		} `json:"packTab2/addScore/index.html"`
		PackTab2AddHousingIndexHtml struct {
			Window struct {
			} `json:"window"`
		} `json:"packTab2/addHousing/index.html"`
		PackTab2MovaImgIndexHtml struct {
			Window struct {
			} `json:"window"`
		} `json:"packTab2/movaImg/index.html"`
		PackTab2VideoWatchIndexHtml struct {
			Window struct {
				NavigationStyle        string `json:"navigationStyle"`
				NavigationBarTextStyle string `json:"navigationBarTextStyle"`
			} `json:"window"`
		} `json:"packTab2/videoWatch/index.html"`
		PackTab2OurTalkIndexHtml struct {
			Window struct {
			} `json:"window"`
		} `json:"packTab2/ourTalk/index.html"`
		PackTab2MyTrendsIndexHtml struct {
			Window struct {
				NavigationStyle string `json:"navigationStyle"`
			} `json:"window"`
		} `json:"packTab2/myTrends/index.html"`
		PackTab2ExCommonIndexHtml struct {
			Window struct {
				NavigationStyle string `json:"navigationStyle"`
			} `json:"window"`
		} `json:"packTab2/exCommon/index.html"`
		PackTab2DjzlPageIndexHtml struct {
			Window struct {
			} `json:"window"`
		} `json:"packTab2/djzlPage/index.html"`
		PackTab2GfzlPageIndexHtml struct {
			Window struct {
				NavigationStyle string `json:"navigationStyle"`
			} `json:"window"`
		} `json:"packTab2/gfzlPage/index.html"`
		PackTab2YhCjWtPageIndexHtml struct {
			Window struct {
			} `json:"window"`
		} `json:"packTab2/yhCjWtPage/index.html"`
		PackTab2ZiLiaoMuBanIndexHtml struct {
			Window struct {
				NavigationStyle string `json:"navigationStyle"`
			} `json:"window"`
		} `json:"packTab2/ziLiaoMuBan/index.html"`
		PackTab2YaohaoPageIndexHtml struct {
			Window struct {
			} `json:"window"`
		} `json:"packTab2/yaohaoPage/index.html"`
		PackTab2TopicAskIndexHtml struct {
			Window struct {
			} `json:"window"`
		} `json:"packTab2/topicAsk/index.html"`
		PackTab2LocateMapIndexHtml struct {
			Window struct {
			} `json:"window"`
		} `json:"packTab2/locateMap/index.html"`
		PackTab2MoreTrendIndexHtml struct {
			Window struct {
			} `json:"window"`
		} `json:"packTab2/moreTrend/index.html"`
		PackTab2GetOneVoteIndexHtml struct {
			Window struct {
				NavigationStyle string `json:"navigationStyle"`
			} `json:"window"`
		} `json:"packTab2/getOneVote/index.html"`
		PackTab2GetAllVoteIndexHtml struct {
			Window struct {
				NavigationStyle string `json:"navigationStyle"`
			} `json:"window"`
		} `json:"packTab2/getAllVote/index.html"`
		PackTab2SipanDongtaiIndexHtml struct {
			Window struct {
				NavigationStyle string `json:"navigationStyle"`
			} `json:"window"`
		} `json:"packTab2/sipan_dongtai/index.html"`
		PackTab2OneResultIndexHtml struct {
			Window struct {
				NavigationStyle string `json:"navigationStyle"`
			} `json:"window"`
		} `json:"packTab2/oneResult/index.html"`
		PackTab2GfBaiKeIndexHtml struct {
			Window struct {
				NavigationStyle string `json:"navigationStyle"`
			} `json:"window"`
		} `json:"packTab2/gfBaiKe/index.html"`
		PackTab2BaikeDetailIndexHtml struct {
			Window struct {
				NavigationStyle string `json:"navigationStyle"`
			} `json:"window"`
		} `json:"packTab2/baikeDetail/index.html"`
		PackTab2ArtLanmuIndexHtml struct {
			Window struct {
				NavigationStyle string `json:"navigationStyle"`
			} `json:"window"`
		} `json:"packTab2/artLanmu/index.html"`
		PackTab2MoreCustomerIndexHtml struct {
			Window struct {
			} `json:"window"`
		} `json:"packTab2/moreCustomer/index.html"`
		PackTab2ServiceViewIndexHtml struct {
			Window struct {
			} `json:"window"`
		} `json:"packTab2/serviceView/index.html"`
		PackTab2HouseStockIndexHtml struct {
			Window struct {
			} `json:"window"`
		} `json:"packTab2/houseStock/index.html"`
		PackTab2BindBuildIndexHtml struct {
			Window struct {
			} `json:"window"`
		} `json:"packTab2/bindBuild/index.html"`
		PackTab2ReportIndexHtml struct {
			Window struct {
			} `json:"window"`
		} `json:"packTab2/report/index.html"`
		PackTab2ActivityLAIndexHtml struct {
			Window struct {
			} `json:"window"`
		} `json:"packTab2/activityLA/index.html"`
		PackTab2FeedBackIndexHtml struct {
			Window struct {
				NavigationStyle string `json:"navigationStyle"`
			} `json:"window"`
		} `json:"packTab2/feedBack/index.html"`
		PackTab2XjShareIndexHtml struct {
			Window struct {
			} `json:"window"`
		} `json:"packTab2/xjShare/index.html"`
		PackTab2OwnerApplyIndexHtml struct {
			Window struct {
				NavigationStyle string `json:"navigationStyle"`
			} `json:"window"`
		} `json:"packTab2/ownerApply/index.html"`
		PackTab2ApplyMesIndexHtml struct {
			Window struct {
			} `json:"window"`
		} `json:"packTab2/applyMes/index.html"`
		PackTab2ToArticlesIndexHtml struct {
			Window struct {
			} `json:"window"`
		} `json:"packTab2/toArticles/index.html"`
		PackTab2LiveShowIndexHtml struct {
			Window struct {
			} `json:"window"`
		} `json:"packTab2/liveShow/index.html"`
		PackTab2ChatUsersIndexHtml struct {
			Window struct {
				NavigationStyle string `json:"navigationStyle"`
			} `json:"window"`
		} `json:"packTab2/chatUsers/index.html"`
		PackTab2PjListIndexHtml struct {
			Window struct {
				NavigationStyle string `json:"navigationStyle"`
			} `json:"window"`
		} `json:"packTab2/pjList/index.html"`
		PackTab2ShareIndexHtml struct {
			Window struct {
				NavigationStyle string `json:"navigationStyle"`
			} `json:"window"`
		} `json:"packTab2/share/index.html"`
		PackTab2DatareportingIndexHtml struct {
			Window struct {
				NavigationStyle string `json:"navigationStyle"`
			} `json:"window"`
		} `json:"packTab2/Datareporting/index.html"`
		PackTab2PoiIndexHtml struct {
			Window struct {
				NavigationStyle string `json:"navigationStyle"`
				DisableScroll   bool   `json:"disableScroll"`
			} `json:"window"`
		} `json:"packTab2/poi/index.html"`
		PackTab2ValueMapIndexHtml struct {
			Window struct {
				NavigationStyle string `json:"navigationStyle"`
			} `json:"window"`
		} `json:"packTab2/valueMap/index.html"`
		PackTab2LiveIndexHtml struct {
			Window struct {
				NavigationStyle string `json:"navigationStyle"`
			} `json:"window"`
		} `json:"packTab2/Live/index.html"`
		PackTab2ZxListIndexHtml struct {
			Window struct {
				NavigationStyle string `json:"navigationStyle"`
			} `json:"window"`
		} `json:"packTab2/zxList/index.html"`
		PackTab2HousingIndexHtml struct {
			Window struct {
				NavigationStyle string `json:"navigationStyle"`
			} `json:"window"`
		} `json:"packTab2/housing/index.html"`
		PackTab2ProclamationsIndexHtml struct {
			Window struct {
				NavigationStyle       string `json:"navigationStyle"`
				EnablePullDownRefresh bool   `json:"enablePullDownRefresh"`
			} `json:"window"`
		} `json:"packTab2/Proclamations/index.html"`
		PackTab2BankIndexHtml struct {
			Window struct {
				NavigationStyle string `json:"navigationStyle"`
			} `json:"window"`
		} `json:"packTab2/bank/index.html"`
		PackTab2BaseInfoIndexHtml struct {
			Window struct {
				NavigationStyle string `json:"navigationStyle"`
			} `json:"window"`
		} `json:"packTab2/baseInfo/index.html"`
		PackTab2QVoteIndexHtml struct {
			Window struct {
				NavigationStyle string `json:"navigationStyle"`
			} `json:"window"`
		} `json:"packTab2/qVote/index.html"`
		PackTab2SunsIndexHtml struct {
			Window struct {
				NavigationStyle string `json:"navigationStyle"`
			} `json:"window"`
		} `json:"packTab2/suns/index.html"`
		PackTab2SqvipIndexHtml struct {
			Window struct {
				NavigationStyle string `json:"navigationStyle"`
			} `json:"window"`
		} `json:"packTab2/sqvip/index.html"`
		PackTab2InfoPKIndexHtml struct {
			Window struct {
				NavigationStyle string `json:"navigationStyle"`
			} `json:"window"`
		} `json:"packTab2/InfoPK/index.html"`
		PackTab2HouseSelectorIndexHtml struct {
			Window struct {
				NavigationStyle string `json:"navigationStyle"`
			} `json:"window"`
		} `json:"packTab2/houseSelector/index.html"`
		PackTab2TalklistIndexHtml struct {
			Window struct {
				NavigationStyle string `json:"navigationStyle"`
			} `json:"window"`
		} `json:"packTab2/talklist/index.html"`
		PackTab2MyMoneyIndexHtml struct {
			Window struct {
				NavigationStyle string `json:"navigationStyle"`
			} `json:"window"`
		} `json:"packTab2/myMoney/index.html"`
		PackTab2EnvelopePoolIndexHtml struct {
			Window struct {
				NavigationStyle string `json:"navigationStyle"`
			} `json:"window"`
		} `json:"packTab2/envelopePool/index.html"`
		PackTab2CategoryIndexHtml struct {
			Window struct {
				NavigationStyle string `json:"navigationStyle"`
			} `json:"window"`
		} `json:"packTab2/category/index.html"`
		PackTab2TransactionContractIndexHtml struct {
			Window struct {
				NavigationStyle string `json:"navigationStyle"`
			} `json:"window"`
		} `json:"packTab2/transactionContract/index.html"`
		PackTab3EvaluationListIndexHtml struct {
			Window struct {
				NavigationStyle       string `json:"navigationStyle"`
				EnablePullDownRefresh bool   `json:"enablePullDownRefresh"`
				DisableScroll         bool   `json:"disableScroll"`
			} `json:"window"`
		} `json:"packTab3/evaluationList/index.html"`
		PackTab3AsianGamesRedPacketIndexHtml struct {
			Window struct {
				NavigationStyle       string `json:"navigationStyle"`
				EnablePullDownRefresh bool   `json:"enablePullDownRefresh"`
			} `json:"window"`
		} `json:"packTab3/asianGames_redPacket/index.html"`
		PackTab3CouponAgreementIndexHtml struct {
			Window struct {
				NavigationStyle       string `json:"navigationStyle"`
				EnablePullDownRefresh bool   `json:"enablePullDownRefresh"`
			} `json:"window"`
		} `json:"packTab3/couponAgreement/index.html"`
		PackTab3NotificationManagerIndexHtml struct {
			Window struct {
				NavigationStyle       string `json:"navigationStyle"`
				EnablePullDownRefresh bool   `json:"enablePullDownRefresh"`
				DisableScroll         bool   `json:"disableScroll"`
			} `json:"window"`
		} `json:"packTab3/notificationManager/index.html"`
		PackTab3ShopBrandIndexHtml struct {
			Window struct {
				NavigationStyle string `json:"navigationStyle"`
			} `json:"window"`
		} `json:"packTab3/shopBrand/index.html"`
		PackTab3HouseTagsIndexHtml struct {
			Window struct {
				NavigationStyle string `json:"navigationStyle"`
			} `json:"window"`
		} `json:"packTab3/houseTags/index.html"`
		PackTab3QwcSouIndexHtml struct {
			Window struct {
				NavigationStyle string `json:"navigationStyle"`
			} `json:"window"`
		} `json:"packTab3/qwcSou/index.html"`
		PackTab3ImgMapIndexHtml struct {
			Window struct {
				NavigationStyle string `json:"navigationStyle"`
			} `json:"window"`
		} `json:"packTab3/imgMap/index.html"`
		PackTab3FloorPlansShopsIndexHtml struct {
			Window struct {
				NavigationStyle string `json:"navigationStyle"`
			} `json:"window"`
		} `json:"packTab3/floorPlansShops/index.html"`
		PackTab3OneStoreIndexHtml struct {
			Window struct {
				NavigationStyle string `json:"navigationStyle"`
			} `json:"window"`
		} `json:"packTab3/oneStore/index.html"`
		PackTab3MyRedbacketIndexHtml struct {
			Window struct {
				NavigationBarTextStyle string `json:"navigationBarTextStyle"`
				NavigationStyle        string `json:"navigationStyle"`
			} `json:"window"`
		} `json:"packTab3/myRedbacket/index.html"`
		PackTab3RedPacketReviewIndexHtml struct {
			Window struct {
				NavigationBarTextStyle string `json:"navigationBarTextStyle"`
				NavigationStyle        string `json:"navigationStyle"`
			} `json:"window"`
		} `json:"packTab3/redPacketReview/index.html"`
		PackTab3LotteryRedEnvelopeIndexHtml struct {
			Window struct {
				NavigationStyle string `json:"navigationStyle"`
			} `json:"window"`
		} `json:"packTab3/lotteryRedEnvelope/index.html"`
		PackTab3MyTopicIndexHtml struct {
			Window struct {
				NavigationStyle string `json:"navigationStyle"`
			} `json:"window"`
		} `json:"packTab3/myTopic/index.html"`
		PackTab3MyExamineIndexHtml struct {
			Window struct {
				NavigationStyle string `json:"navigationStyle"`
			} `json:"window"`
		} `json:"packTab3/myExamine/index.html"`
		PackTab3CardsIndexHtml struct {
			Window struct {
				NavigationStyle string `json:"navigationStyle"`
			} `json:"window"`
		} `json:"packTab3/cards/index.html"`
		PackTab3DailyTranspondIndexHtml struct {
			Window struct {
				NavigationStyle       string `json:"navigationStyle"`
				EnablePullDownRefresh bool   `json:"enablePullDownRefresh"`
			} `json:"window"`
		} `json:"packTab3/dailyTranspond/index.html"`
		PackTab3SecHandHouseIndexHtml struct {
			Window struct {
				NavigationStyle       string `json:"navigationStyle"`
				EnablePullDownRefresh bool   `json:"enablePullDownRefresh"`
			} `json:"window"`
		} `json:"packTab3/secHandHouse/index.html"`
		PackTab3EstateHistoryDealDataIndexHtml struct {
			Window struct {
				NavigationStyle       string `json:"navigationStyle"`
				EnablePullDownRefresh bool   `json:"enablePullDownRefresh"`
				OnReachBottomDistance int    `json:"onReachBottomDistance"`
			} `json:"window"`
		} `json:"packTab3/estateHistoryDealData/index.html"`
		PackTab3HouseEvoluationIndexHtml struct {
			Window struct {
				NavigationStyle       string `json:"navigationStyle"`
				EnablePullDownRefresh bool   `json:"enablePullDownRefresh"`
			} `json:"window"`
		} `json:"packTab3/houseEvoluation/index.html"`
		PackTab3XiaojiSuperSaleIndexHtml struct {
			Window struct {
				NavigationStyle       string `json:"navigationStyle"`
				EnablePullDownRefresh bool   `json:"enablePullDownRefresh"`
			} `json:"window"`
		} `json:"packTab3/xiaojiSuperSale/index.html"`
		PackTab3TimeLimitedTaskIndexHtml struct {
			Window struct {
				NavigationStyle       string `json:"navigationStyle"`
				EnablePullDownRefresh bool   `json:"enablePullDownRefresh"`
			} `json:"window"`
		} `json:"packTab3/timeLimitedTask/index.html"`
		PackTab3ShoppingMallIndexHtml struct {
			Window struct {
				NavigationStyle       string `json:"navigationStyle"`
				EnablePullDownRefresh bool   `json:"enablePullDownRefresh"`
			} `json:"window"`
		} `json:"packTab3/shoppingMall/index.html"`
		PackTab3BrandPavilionItemIndexHtml struct {
			Window struct {
				NavigationStyle string `json:"navigationStyle"`
			} `json:"window"`
		} `json:"packTab3/brandPavilion_item/index.html"`
		PackTab3BrandPavilionIndexHtml struct {
			Window struct {
				NavigationStyle string `json:"navigationStyle"`
			} `json:"window"`
		} `json:"packTab3/brandPavilion/index.html"`
		PackTab3ConsultantShareIndexHtml struct {
			Window struct {
				NavigationStyle       string `json:"navigationStyle"`
				OnReachBottomDistance int    `json:"onReachBottomDistance"`
			} `json:"window"`
		} `json:"packTab3/consultantShare/index.html"`
		PackTab3HomepageMapIndexHtml struct {
			Window struct {
				NavigationStyle string `json:"navigationStyle"`
			} `json:"window"`
		} `json:"packTab3/homepageMap/index.html"`
		PackTab3RedLotteryIndexHtml struct {
			Window struct {
				NavigationStyle string `json:"navigationStyle"`
			} `json:"window"`
		} `json:"packTab3/redLottery/index.html"`
		PackTab3MyCustomerIndexHtml struct {
			Window struct {
				NavigationStyle string `json:"navigationStyle"`
			} `json:"window"`
		} `json:"packTab3/myCustomer/index.html"`
		PackTab3RedActiveIndexHtml struct {
			Window struct {
				NavigationStyle string `json:"navigationStyle"`
			} `json:"window"`
		} `json:"packTab3/redActive/index.html"`
		PackTab3ProgrammeIndexHtml struct {
			Window struct {
				NavigationStyle string `json:"navigationStyle"`
			} `json:"window"`
		} `json:"packTab3/programme/index.html"`
		PackTab3ShareHbIndexHtml struct {
			Window struct {
				NavigationStyle string `json:"navigationStyle"`
			} `json:"window"`
		} `json:"packTab3/shareHb/index.html"`
		PackTab3MineIndexHtml struct {
			Window struct {
				NavigationStyle string `json:"navigationStyle"`
			} `json:"window"`
		} `json:"packTab3/Mine/index.html"`
		PackTab3SchemeRecordIndexHtml struct {
			Window struct {
				NavigationStyle string `json:"navigationStyle"`
			} `json:"window"`
		} `json:"packTab3/schemeRecord/index.html"`
		PackTab3SchemeIndexHtml struct {
			Window struct {
				NavigationStyle string `json:"navigationStyle"`
			} `json:"window"`
		} `json:"packTab3/scheme/index.html"`
		PackTab3ShareCustomerIndexHtml struct {
			Window struct {
				NavigationStyle string `json:"navigationStyle"`
			} `json:"window"`
		} `json:"packTab3/shareCustomer/index.html"`
		PackTab3HouseInfoIndexHtml struct {
			Window struct {
				NavigationStyle string `json:"navigationStyle"`
			} `json:"window"`
		} `json:"packTab3/houseInfo/index.html"`
		PackTab3LptopIndexHtml struct {
			Window struct {
				NavigationStyle string `json:"navigationStyle"`
			} `json:"window"`
		} `json:"packTab3/lptop/index.html"`
		PackTab3PlanValueIndexHtml struct {
			Window struct {
				NavigationStyle string `json:"navigationStyle"`
			} `json:"window"`
		} `json:"packTab3/planValue/index.html"`
		PackTab3PlanLocationIndexHtml struct {
			Window struct {
				NavigationStyle string `json:"navigationStyle"`
			} `json:"window"`
		} `json:"packTab3/planLocation/index.html"`
		PackTab3PlanNearIndexHtml struct {
			Window struct {
				NavigationStyle string `json:"navigationStyle"`
			} `json:"window"`
		} `json:"packTab3/planNear/index.html"`
		PackTab3PlanAreaIndexHtml struct {
			Window struct {
				NavigationStyle string `json:"navigationStyle"`
			} `json:"window"`
		} `json:"packTab3/planArea/index.html"`
		PackTab3PlanInfoIndexHtml struct {
			Window struct {
				NavigationStyle string `json:"navigationStyle"`
			} `json:"window"`
		} `json:"packTab3/planInfo/index.html"`
		PackTab3PlanSummaryIndexHtml struct {
			Window struct {
				NavigationStyle string `json:"navigationStyle"`
			} `json:"window"`
		} `json:"packTab3/planSummary/index.html"`
		PackTab3PlanApartIndexHtml struct {
			Window struct {
				NavigationStyle string `json:"navigationStyle"`
			} `json:"window"`
		} `json:"packTab3/planApart/index.html"`
		PackTab3HotmoreIndexHtml struct {
			Window struct {
				NavigationStyle string `json:"navigationStyle"`
			} `json:"window"`
		} `json:"packTab3/hotmore/index.html"`
		PackTab3NewActiveIndexHtml struct {
			Window struct {
				NavigationStyle string `json:"navigationStyle"`
			} `json:"window"`
		} `json:"packTab3/NewActive/index.html"`
		PackTab3TenancyCommonIndexHtml struct {
			Window struct {
				NavigationStyle string `json:"navigationStyle"`
			} `json:"window"`
		} `json:"packTab3/tenancyCommon/index.html"`
		PackTab3ParkingpacesIndexHtml struct {
			Window struct {
				NavigationStyle string `json:"navigationStyle"`
			} `json:"window"`
		} `json:"packTab3/parkingpaces/index.html"`
		PackTab3DrivewayClientIndexHtml struct {
			Window struct {
				NavigationStyle string `json:"navigationStyle"`
			} `json:"window"`
		} `json:"packTab3/drivewayClient/index.html"`
		PackTab3BuildingCommentIndexHtml struct {
			Window struct {
				NavigationStyle string `json:"navigationStyle"`
			} `json:"window"`
		} `json:"packTab3/buildingComment/index.html"`
		PackTab3ChangeSellMapIndexHtml struct {
			Window struct {
				NavigationStyle string `json:"navigationStyle"`
			} `json:"window"`
		} `json:"packTab3/changeSellMap/index.html"`
		PackTab3ChoseRealEstateIndexHtml struct {
			Window struct {
				NavigationStyle string `json:"navigationStyle"`
			} `json:"window"`
		} `json:"packTab3/choseRealEstate/index.html"`
		PackTab3PrepaymentIndexHtml struct {
			Window struct {
				NavigationStyle       string `json:"navigationStyle"`
				EnablePullDownRefresh bool   `json:"enablePullDownRefresh"`
				DisableScroll         bool   `json:"disableScroll"`
			} `json:"window"`
		} `json:"packTab3/prepayment/index.html"`
		PackTab3PrepayDetailIndexHtml struct {
			Window struct {
				NavigationStyle string `json:"navigationStyle"`
			} `json:"window"`
		} `json:"packTab3/prepayDetail/index.html"`
		PackTab3RulePageIndexHtml struct {
			Window struct {
				NavigationStyle string `json:"navigationStyle"`
			} `json:"window"`
		} `json:"packTab3/rulePage/index.html"`
		PackTab3ShopCollectionIndexHtml struct {
			Window struct {
				NavigationStyle string `json:"navigationStyle"`
			} `json:"window"`
		} `json:"packTab3/shopCollection/index.html"`
		PackTab3AreaShopMapIndexHtml struct {
			Window struct {
				NavigationStyle string `json:"navigationStyle"`
			} `json:"window"`
		} `json:"packTab3/areaShopMap/index.html"`
		PackTab3SouQunIndexHtml struct {
			Window struct {
				NavigationStyle string `json:"navigationStyle"`
			} `json:"window"`
		} `json:"packTab3/souQun/index.html"`
		PackTab3HotMapIndexHtml struct {
			Window struct {
				NavigationStyle string `json:"navigationStyle"`
			} `json:"window"`
		} `json:"packTab3/hotMap/index.html"`
		PackTab3NowPublicityIndexHtml struct {
			Window struct {
				NavigationStyle string `json:"navigationStyle"`
			} `json:"window"`
		} `json:"packTab3/nowPublicity/index.html"`
		PackTab3NowRegisterIndexHtml struct {
			Window struct {
				NavigationStyle string `json:"navigationStyle"`
			} `json:"window"`
		} `json:"packTab3/nowRegister/index.html"`
		PackTab3LotteryQueryIndexHtml struct {
			Window struct {
				NavigationStyle string `json:"navigationStyle"`
			} `json:"window"`
		} `json:"packTab3/lotteryQuery/index.html"`
		PackTab3NewApartmentsIndexHtml struct {
			Window struct {
				NavigationStyle string `json:"navigationStyle"`
			} `json:"window"`
		} `json:"packTab3/newApartments/index.html"`
		PackTab3BuildingRankingIndexHtml struct {
			Window struct {
				NavigationStyle string `json:"navigationStyle"`
			} `json:"window"`
		} `json:"packTab3/buildingRanking/index.html"`
		PackTab3ExpectedOpeningIndexHtml struct {
			Window struct {
				NavigationStyle string `json:"navigationStyle"`
			} `json:"window"`
		} `json:"packTab3/expectedOpening/index.html"`
		PackTab3DivideQuataIndexHtml struct {
			Window struct {
				NavigationStyle string `json:"navigationStyle"`
			} `json:"window"`
		} `json:"packTab3/divideQuata/index.html"`
		PackTab3ChatHistoryIndexHtml struct {
			Window struct {
				NavigationStyle string `json:"navigationStyle"`
			} `json:"window"`
		} `json:"packTab3/chatHistory/index.html"`
		PackTab3ChickPreferenceIndexHtml struct {
			Window struct {
				NavigationStyle string `json:"navigationStyle"`
			} `json:"window"`
		} `json:"packTab3/chickPreference/index.html"`
		PackTab3SuperSaleReviceIndexHtml struct {
			Window struct {
				NavigationStyle string `json:"navigationStyle"`
			} `json:"window"`
		} `json:"packTab3/superSaleRevice/index.html"`
		PackTab3HistoryDataIndexHtml struct {
			Window struct {
				NavigationBarTitleText string `json:"navigationBarTitleText"`
				NavigationStyle        string `json:"navigationStyle"`
				EnablePullDownRefresh  bool   `json:"enablePullDownRefresh"`
				DisableScroll          bool   `json:"disableScroll"`
				Renderer               string `json:"renderer"`
				ComponentFramework     string `json:"componentFramework"`
			} `json:"window"`
		} `json:"packTab3/HistoryData/index.html"`
		PackTab3QwcIndexHtml struct {
			Window struct {
				NavigationStyle string `json:"navigationStyle"`
			} `json:"window"`
		} `json:"packTab3/qwc/index.html"`
		PackTab3ShhandSearchIndexHtml struct {
			Window struct {
				NavigationStyle       string `json:"navigationStyle"`
				EnablePullDownRefresh bool   `json:"enablePullDownRefresh"`
			} `json:"window"`
		} `json:"packTab3/shhandSearch/index.html"`
		PackTab3SuggestionIndexHtml struct {
			Window struct {
				NavigationStyle string `json:"navigationStyle"`
			} `json:"window"`
		} `json:"packTab3/suggestion/index.html"`
		PackTab3FeedBackRecordIndexHtml struct {
			Window struct {
				NavigationStyle string `json:"navigationStyle"`
			} `json:"window"`
		} `json:"packTab3/feedBackRecord/index.html"`
		PackTab3FeedbackSubmitIndexHtml struct {
			Window struct {
				NavigationStyle       string `json:"navigationStyle"`
				EnablePullDownRefresh bool   `json:"enablePullDownRefresh"`
			} `json:"window"`
		} `json:"packTab3/feedbackSubmit/index.html"`
		PackTab3SuggestionsManageIndexHtml struct {
			Window struct {
				NavigationStyle       string `json:"navigationStyle"`
				EnablePullDownRefresh bool   `json:"enablePullDownRefresh"`
			} `json:"window"`
		} `json:"packTab3/suggestionsManage/index.html"`
		PackTab2PluginWx2B03C6E691Cd7370PagesLivePlayerPluginHtml struct {
			Window struct {
				UsingComponents struct {
					PageLivePlayer              string `json:"page-live-player"`
					PageLiveReplay              string `json:"page-live-replay"`
					PageLiveRecord              string `json:"page-live-record"`
					ComponentProfileCard        string `json:"component-profile-card"`
					ComponentPersonOperation    string `json:"component-person-operation"`
					ComponentComments           string `json:"component-comments"`
					ComponentStoreList          string `json:"component-store-list"`
					ComponentProfileModal       string `json:"component-profile-modal"`
					ComponentBarrageList        string `json:"component-barrage-list"`
					ComponentGoodsPush          string `json:"component-goods-push"`
					ComponentCouponPush         string `json:"component-coupon-push"`
					ComponentPushComment        string `json:"component-push-comment"`
					ComponentEndBlock           string `json:"component-end-block"`
					ComponentSubscribeCard      string `json:"component-subscribe-card"`
					ComponentForbidUserList     string `json:"component-forbid-user-list"`
					ComponentLottery            string `json:"component-lottery"`
					ComponentLotteryOper        string `json:"component-lottery-oper"`
					ComponentLotteryFillPhone   string `json:"component-lottery-fill-phone"`
					ComponentCommentActionSheet string `json:"component-comment-action-sheet"`
					ComponentTopPlayer          string `json:"component-top-player"`
					ComponentOperatePaster      string `json:"component-operate-paster"`
					ComponentAttentionGuide     string `json:"component-attention-guide"`
					ComponentSettingMore        string `json:"component-setting-more"`
					ComponentMicTag             string `json:"component-mic-tag"`
					ComponentMicPanel           string `json:"component-mic-panel"`
					ComponentAutoReplyPanel     string `json:"component-auto-reply-panel"`
					ComponentMicSystemMsg       string `json:"component-mic-system-msg"`
					ComponentActivityNudge      string `json:"component-activity-nudge"`
					ComponentSharePanel         string `json:"component-share-panel"`
					ComponentNoticeToptips      string `json:"component-notice-toptips"`
					ComponentAdvanceDetail      string `json:"component-advance-detail"`
					ComponentAdvance            string `json:"component-advance"`
					ComponentAutoReply          string `json:"component-auto-reply"`
					ComponentHelpListEntr       string `json:"component-help-list-entr"`
					ComponentHelpPanel          string `json:"component-help-panel"`
					ComponentHelpRulePanel      string `json:"component-help-rule-panel"`
					MpIcon                      string `json:"mp-icon"`
				} `json:"usingComponents"`
				NavigationStyle        string `json:"navigationStyle"`
				BackgroundColor        string `json:"backgroundColor"`
				DisableScroll          bool   `json:"disableScroll"`
				Style                  string `json:"style"`
				NavigationBarTextStyle string `json:"navigationBarTextStyle"`
				SinglePage             struct {
					NavigationBarBackgroundColor string `json:"navigationBarBackgroundColor"`
					NavigationBarBackgroundAlpha int    `json:"navigationBarBackgroundAlpha"`
					NavigationBarTextStyle       string `json:"navigationBarTextStyle"`
				} `json:"singlePage"`
				PageOrientation      string `json:"pageOrientation"`
				HandleWebviewPreload string `json:"handleWebviewPreload"`
				Warning              string `json:"__warning__"`
			} `json:"window"`
		} `json:"packTab2/__plugin__/wx2b03c6e691cd7370/pages/live-player-plugin.html"`
		PackTab2PluginWx2B03C6E691Cd7370WxliveComponentsExtPlayerAddressPreviewAddressPreviewHtml struct {
			Window struct {
				UsingComponents struct {
					MpMsg           string `json:"mp-msg"`
					MpNavigationBar string `json:"mp-navigation-bar"`
				} `json:"usingComponents"`
				NavigationStyle string `json:"navigationStyle"`
				DisableScroll   bool   `json:"disableScroll"`
				Style           string `json:"style"`
			} `json:"window"`
		} `json:"packTab2/__plugin__/wx2b03c6e691cd7370/wxlive-components/ext-player/address-preview/address-preview.html"`
		PackTab2PluginWx2B03C6E691Cd7370WxliveComponentsExtPlayerComplaintCommentComplaintCommentHtml struct {
			Window struct {
				UsingComponents struct {
					MpMsg           string `json:"mp-msg"`
					MpCheckboxGroup string `json:"mp-checkbox-group"`
					MpCheckbox      string `json:"mp-checkbox"`
					MpCells         string `json:"mp-cells"`
					MpNavigationBar string `json:"mp-navigation-bar"`
				} `json:"usingComponents"`
				NavigationStyle string `json:"navigationStyle"`
				DisableScroll   bool   `json:"disableScroll"`
				Style           string `json:"style"`
			} `json:"window"`
		} `json:"packTab2/__plugin__/wx2b03c6e691cd7370/wxlive-components/ext-player/complaint-comment/complaint-comment.html"`
		PackTab2PluginWx2B03C6E691Cd7370WxliveComponentsExtPlayerComplaintRoomComplaintRoomHtml struct {
			Window struct {
				UsingComponents struct {
					MpMsg           string `json:"mp-msg"`
					MpCheckboxGroup string `json:"mp-checkbox-group"`
					MpCheckbox      string `json:"mp-checkbox"`
					MpCells         string `json:"mp-cells"`
					MpNavigationBar string `json:"mp-navigation-bar"`
				} `json:"usingComponents"`
				NavigationStyle string `json:"navigationStyle"`
				DisableScroll   bool   `json:"disableScroll"`
				Style           string `json:"style"`
			} `json:"window"`
		} `json:"packTab2/__plugin__/wx2b03c6e691cd7370/wxlive-components/ext-player/complaint-room/complaint-room.html"`
	} `json:"page"`
	Permission struct {
		ScopeUserLocation struct {
			Desc string `json:"desc"`
		} `json:"scope.userLocation"`
	} `json:"permission"`
	Global struct {
		Window struct {
			BackgroundTextStyle          string `json:"backgroundTextStyle"`
			BackgroundColor              string `json:"backgroundColor"`
			NavigationBarBackgroundColor string `json:"navigationBarBackgroundColor"`
			NavigationBarTitleText       string `json:"navigationBarTitleText"`
			NavigationBarTextStyle       string `json:"navigationBarTextStyle"`
			EnablePullDownRefresh        bool   `json:"enablePullDownRefresh"`
			OnReachBottomDistance        int    `json:"onReachBottomDistance"`
		} `json:"window"`
	} `json:"global"`
	Plugins struct {
		LivePlayerPlugin struct {
			Version    string `json:"version"`
			Provider   string `json:"provider"`
			Subpackage string `json:"subpackage"`
		} `json:"live-player-plugin"`
	} `json:"plugins"`
	TabBar struct {
		Color           string `json:"color"`
		SelectedColor   string `json:"selectedColor"`
		BackgroundColor string `json:"backgroundColor"`
		BorderStyle     string `json:"borderStyle"`
		List            []struct {
			PagePath         string `json:"pagePath"`
			Text             string `json:"text"`
			IconData         string `json:"iconData"`
			SelectedIconData string `json:"selectedIconData"`
		} `json:"list"`
		Position string `json:"position"`
	} `json:"tabBar"`
	RendererOptions struct {
		Skyline struct {
			DefaultDisplayBlock bool `json:"defaultDisplayBlock"`
		} `json:"skyline"`
	} `json:"rendererOptions"`
	Renderer struct {
		AllUsed []string `json:"allUsed"`
		Default string   `json:"default"`
	} `json:"renderer"`
	RequiredPrivateInfos []string `json:"requiredPrivateInfos"`
	UsePrivacyCheck      bool     `json:"__usePrivacyCheck__"`
}
