// app.js
App({
  onLaunch() {
    // 展示本地存储能力
    const logs = wx.getStorageSync('logs') || []
    logs.unshift(Date.now())
    wx.setStorageSync('logs', logs)

  },
  globalData: {
    userInfo: null,

    isLogin: false,
    LoginToken: '',

    request_header: 'https://mini.81jcpd.cn',
    request_header2: 'https://calc.81jcpd.cn',
    question_reflash: true,
    vote_reflash: true,
    vote_detail_refalsh: true,
    votePublish: false,
    votePublishIdx: -1,
    delVoteFlag: false,
    delVoteIdx: -1,
    loginToSW: '',

    uc: '',
    uc_sign: '',

    pk_s: {
      N: -1,
      G: '',
      Y: '',
      H: 0,
      ct_u: [],
      ct_y: []
    }

  }
});

















