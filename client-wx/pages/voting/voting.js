const app = getApp();

Page({
  data: {
    islogin: false,
    loading: false,
    vid: '',
    title: '',
    optionNames: [],
    imc: 0,
    status: 0,
    optionL: ['A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K'],
    selectedIndexs: [],
    optionTypes: [],
    vvvoteLoading: false
  },

  onLoad(options) {
    this.setData({ vid: options['vid'] });
  },

  onShow() {
    if (!app.globalData.isLogin) {
      this.setData({ islogin: false });
      return;
    }
    this.setData({ islogin: true });
    this.loadData();
  },

  toLogin(e) {
    if (!app.globalData.isLogin) {
      let topath = '/pages/voting/voting?vid=' + this.data.vid;
      app.globalData.loginToSW = encodeURIComponent(topath);
      wx.switchTab({
        url: '/pages/index/index',
      });
    }
  },

  loadData() {
    if (this.data.loading) {
      return;
    }
    this.setData({ loading: true });

    var that = this;

    wx.showLoading({
      title: '加载中',
      mask: true
    });

    wx.request({
      url: app.globalData.request_header + '/vote/get-detail',
      method: 'GET',
      header: { "X-auth": app.globalData.LoginToken },
      data: { "id": that.data.vid },
      success(res) {
        if (res.data.code !== 200) {
          wx.showToast({
            title: res.data.msg,
            icon: 'none',
            duration: 1500
          })
          return
        }

        that.setData({
          optionNames: res.data.data.options,
          title: res.data.data.title,
          status: res.data.data.status,
          imc: res.data.data.is_multi_choice
        });

        let types = [];
        let selects = [];
        for (let i = 0; i < that.data.optionNames.length; i++) {
          types.push(0);
          selects.push(0);
        }
        that.setData({ optionTypes: types, selectedIndexs: selects });
      },
      fail(err) {
        wx.showToast({
          title: '服务异常，请稍后再试',
          icon: 'none',
          duration: 1000
        })
      },
      complete() {
        that.setData({ loading: false });
      }
    });

    wx.hideLoading();
  },

  selectOption(e) {
    let idx = e.currentTarget.dataset.idx;
    let types = this.data.optionTypes;
    let selectedIdxs = this.data.selectedIndexs;

    if (selectedIdxs[idx] === 0) {
      if (this.data.imc === 0) {
        cleanIntArr(types);
        cleanIntArr(selectedIdxs);
        types[idx] = 1;
        selectedIdxs[idx] = 1;
      } else {
        types[idx] ^= 1;
        selectedIdxs[idx] ^= 1;
      }
    } else {
      selectedIdxs[idx] = 0;
      types[idx] = 0;
    }

    this.setData({ optionTypes: types, selectedIndexs: selectedIdxs });
  },


  doVote(e) {
    //  检查发布状态
    if (this.data.status === 0) {
      wx.showToast({
        title: '该投票未发布或已停止发布',
        icon: 'none',
        duration: 1000
      });
      return
    }
    //  检查是否至少选择了一个选项
    let selectIdx = this.data.selectedIndexs;
    let choose_num = 0;
    for (let i = 0; i < selectIdx.length; ++i) {
      if (selectIdx[i] !== 0) {
        choose_num++;
      }
    }
    if (choose_num < 1) {
      wx.showToast({
        title: '至少选择一个选项',
        icon: 'none',
        duration: 1000
      });
      return
    }
    //  检查选择数是否满足单选
    if (this.data.imc === 0 && choose_num !== 1) {
      wx.showToast({
        title: '你只能选择一项进行投票',
        icon: 'none',
        duration: 1000
      });
      return
    }
    //  前置条件满足，计算投票的msg
    let msg = calc_msg(selectIdx);
    //  投票
    wx.showLoading({
      title: '正在火速投票',
      mask: true
    });
    vvvote(this.data.vid, msg);
    wx.hideLoading();
  }

});

function cleanIntArr(arr) {
  for (let i = 0; i < arr.length; i++) {
    arr[i] = 0;
  }
}

// --------

function power(base, exponent) {
  base = BigInt(base);
  exponent = BigInt(exponent);
  let result = 1n;
  let currentPower = base;
  while (exponent > 0) {
    if (exponent % 2n === 1n) {
      result *= currentPower;
    }
    currentPower *= currentPower;
    exponent /= 2n;
  }
  return result;
}

function ct(g_x, y_x, h_x) {
  let G = BigInt(app.globalData.pk_s.G);
  let Y = BigInt(app.globalData.pk_s.Y);
  let H = BigInt(app.globalData.pk_s.H);

  let gx = BigInt(g_x);
  let yx = BigInt(y_x);
  let hx = BigInt(h_x);

  const capacity = 2;
  const ct_x = new Array(capacity);

  ct_x[0] = (power(G, gx)).toString();
  ct_x[1] = (power(Y, yx) * power(H, hx)).toString();

  return ct_x;
}

function num_mul_two_tuple(tuple, num) {
  let num_big = BigInt(num);
  let t1 = BigInt(tuple[0]);
  let t2 = BigInt(tuple[1]);

  const r = [];
  r[0] = power(t1, num_big);
  r[1] = power(t2, num_big);

  return r;
}

function merge_ct(ct1, ct2, ct3) {
  const r = [];
  r[0] = BigInt(ct1[0]) * BigInt(ct2[0]) * BigInt(ct3[0]);
  r[1] = BigInt(ct1[1]) * BigInt(ct2[1]) * BigInt(ct3[1]);
  return r;
}

function unblind(f, a) {
  let beta_ie = -1;
  let F = BigInt(f);
  for (let i = 1; i < app.globalData.pk_s.N; i++) {
    let ii = BigInt(i);
    if (power(app.globalData.pk_s.G, ii).toString() === F.toString()) {
      beta_ie = i;
    }
  }

  if (beta_ie === -1) {
    return [false, '']
  }

  let bap = (beta_ie * a) % app.globalData.pk_s.N;
  return [true, power(app.globalData.pk_s.G, BigInt(bap)).toString()]
}

function getRandomIntInRange(p) {
  const min = 3;
  const max = p;
  // Math.random() 生成 [0, 1) 之间的随机数
  // Math.floor 将其调整为 [min, max) 之间的整数
  const randomInt = Math.floor(Math.random() * (max - min) + min);
  return randomInt;
}

function calc_msg(selectIndexs) {
  let r = 0
  for (let i = 0; i < selectIndexs.length; ++i) {
    if (selectIndexs[i] !== 0) {
      r += (1 << i);
    }
  }
  return r
}

function calc_ct_beta(a, b, r, p, m) {
  let am_bp = a * m + b * p;
  let ct1 = ct(r, r, am_bp);
  let ct2 = num_mul_two_tuple(app.globalData.pk_s.ct_u, a);
  let ct3 = num_mul_two_tuple(app.globalData.pk_s.ct_y, a * r);
  let ct_beta = merge_ct(ct1, ct2, ct3);
  return ct_beta
}

function verify(uc_, hsign_, vid_, msg_, r_, token_) {
  wx.request({
    url: app.globalData.request_header + '/auth/act/verify',
    method: 'POST',
    data: {
      msg: msg_,
      r: r_,
      uc: uc_,
      vid: parseInt(vid_, 10),
      token: token_,
      h_sign: hsign_
    },
    success(res) {
      if (res.data.code !== 200) {
        wx.showToast({
          title: res.data.msg,
          icon: 'none',
          duration: 1500
        });
        return
      }
      wx.showToast({
        title: res.data.data,
        icon: 'none',
        duration: 2000
      });
    },
    fail(err) {
      wx.showToast({
        title: '服务异常，请稍后再试',
        icon: 'none',
        duration: 1000
      });
    }
  });
}

function vvvote(vid, msg) {
  if (this.data.vvvoteLoading) {
    return
  }
  this.setData({ vvvoteLoading: true });
  let a = getRandomIntInRange(app.globalData.pk_s.N / 2);
  let b = getRandomIntInRange(app.globalData.pk_s.N / 2);
  let r = getRandomIntInRange(app.globalData.pk_s.N / 2);

  let ct_beta = calc_ct_beta(a, b, r, app.globalData.pk_s.N, msg);

  wx.request({
    url: app.globalData.request_header + '/auth/act/sign',
    method: 'POST',
    data: {
      u: ct_beta[0].toString(),
      e: ct_beta[1].toString(),
      uc: app.globalData.uc,
      uc_sign: app.globalData.uc_sign,
      vid: parseInt(vid, 10)
    },
    success(res) {
      if (res.data.code !== 200) {
        console.log("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa - sign")
        wx.showToast({
          title: res.data.msg,
          icon: 'none',
          duration: 1000
        });
        return
      }

      let F = res.data.data.F;
      let hsign = res.data.data.h_sign;

      if (F === undefined || F === null || hsign === undefined || hsign === null) {
        wx.showToast({
          title: '数据异常',
          icon: 'none',
          duration: 1000
        });
        return
      }

      const [unblind_f, token] = unblind(F, a);
      if (!unblind_f) {
        wx.showToast({
          title: '数据异常,请稍后再试',
          icon: 'none',
          duration: 1000
        });
        return
      }

      verify(app.globalData.uc, hsign, vid, msg, r, token);
    },
    fail(err) {
      wx.showToast({
        title: '服务异常，请稍后再试',
        icon: 'none',
        duration: 1000
      });
    }
  });
  this.setData({ vvvoteLoading: false });
}

