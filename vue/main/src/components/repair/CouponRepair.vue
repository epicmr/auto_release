<template>
    <div v-loading="this.$store.state.loading">
        <el-row>
            <el-input class="input" v-on:input="OnInput" v-model.number="phone" placeholder="请输入用户手机号码" @keyup.enter.native="Query()" clearable></el-input>
            <el-button class="button" type="primary" @click="Query()" accesskey="1">查询</el-button>
            <el-button class="button" type="primary" @click="Submit()" v-if="canSubmit">{{submitText}}</el-button>
        </el-row>
        <div class="infowrapper" v-if="queryed">
            <el-row class="info">
                <el-col :span=4>
                    <span>优惠券码</span>
                </el-col>
                <el-col :span=4>
                    <span>使用状态</span>
                </el-col>
                <el-col :span=4>
                    <span>订单ID</span>
                </el-col>
                <el-col :span=12>
                    <span>名称</span>
                </el-col>
            </el-row>
            <el-row class="info infoheight">
                <form>
                    <el-row name="coupon_list" :key="r.record_id" v-for="r in this.vipRights">
                        <el-col :span=4>
                            <el-checkbox name="coupon" @change="selected" :label="r.record_id" size="medium"></el-checkbox>
                        </el-col>
                        <el-col :span=4>
                            <el-tag :type="stateTypes[r.use_state]" class="tag" effect="dark">{{stateDescs[r.use_state]}}</el-tag>
                        </el-col>
                        <el-col :span=4>
                            <span>{{r.deal_id}}</span>
                        </el-col>
                        <el-col :span=12>
                            <span>{{r.group_name}}</span>
                        </el-col>
                    </el-row>
                </form>
            </el-row>
        </div>
    </div>
</template>

<script>
    export default {
        name: "v-userrepair",
        data() {
            return {
                userId : 0,
                submitText : '提交',
                canSubmit:false,
                vipRights : '',
                phone : '',
                clist : [],
                stateDescs: {
                    '0': '可使用',
                    '1': '已使用',
                    '2': '已失效'
                },
                stateTypes: {
                    '0': 'success',
                    '1': 'danger',
                    '2': 'warning'
                }
            }
        },
        methods: {
            OnInput(val) {
                this.vipRights = ''
                this.userId = 0
            },
            Query() {
                this.$http
                    .post('http://cgi.gstyun.cn/cgi-bin/user/userquerydetail', {"mobile_number":this.phone})
                    .then(response => {
                        if (parseInt(response.data.status, 10) == 0) {
                            this.userId = response.data.user_id
                        }
                    })
                    .catch(err => {
                        console.info("%c [axiso catch error]", "color:orange", err)
                    });
                let payload = {
                    "page_size":"99",
                    "page_no":"1",
                    "user_id":this.userId,
                }
                this.$http
                    .get('http://stg-cgi.gstyun.cn/cgi-bin/coupon/queryusercouponlist', {params:payload})
                    .then(response => {
                        let couponList = response.data.data.coupon_list
                        if (Array.isArray(couponList)) {
                            this.vipRights = couponList
                        }
                    })
                    .catch(err => {
                        console.info("%c [axiso catch error]", "color:orange", err)
                    });
                //setTimeout(this.selected, 200)
            },
            Submit() {
                let typ = ""
                if (this.submitText == "删除优惠券") {
                    typ = "disable"
                }
                else if (this.submitText == "恢复优惠券") {
                    typ = "enable"
                }
                else {
                    console.log("error")
                }
                this.$http
                    .post('http://stg-go-api.gstyun.cn/coupon/handler/couponState', {"coupon_flts":this.clist, "type":typ})
                    .then(response => {
                        if (parseInt(response.data.code, 10) == 0) {
                            this.$message(this.submitText + '成功');
                        }
                    })
                    .catch(err => {
                        console.info("%c [axiso catch error]", "color:orange", err)
                    });

                //setTimeout(this.Query, 200)
            },
            selected() {
                if (document.forms.length == 0) {
                    this.canSubmit = false;
                    return
                }
                let l = document.forms[0].coupon
                this.canSubmit = false;

                let state = -1
                this.clist = []
                for (let i = 0; i < l.length; i++) {
                    if (l[i].checked) {
                        if (state == -1) {
                            state = this.vipRights[i].use_state == 0 ? 0 : 1
                            console.log("state " + state)
                        }
                        let c = {}
                        let thisState = this.vipRights[i].use_state == 0 ? 0 : 1
                        console.log("this state " + thisState)
                        c.user_id = this.userId
                        c.record_id = this.vipRights[i].record_id
                        this.clist.push(c)
                        if (thisState != state) {
                            this.canSubmit = false;
                            console.log(this.canSubmit)
                            return
                        }
                    }
                }

                if (state != -1) {
                    console.log(this.clist)
                    this.canSubmit = true;
                    this.submitText = (state == 0) ? "删除优惠券":"恢复优惠券"
                    console.log(this.canSubmit)
                }
                return
            }
        },
        computed: {
            queryed : {
                get() {
                    if (this.userId > 0) {
                        return true
                    }
                    return false
                }
            }
        }
    }
</script>

<style scoped>
.input {
    margin: 5px;
    width: 200px;
    border-radius: 4px;
    font-size: 14px;
}
.button {
    margin: 5px;
    border-radius: 4px;
}
.radiogroup
{
    margin-top: 10px;
    margin-bottom: 15px;
}
.infowrapper
{
    padding: 5px;
    border: 1px solid #D7D7D7;
    border-radius: 5px;
    background-color: #E4E4E4;
    font-size: 14px;
    margin: 5px;
}
.info
{
    padding: 5px;
    border: 1px solid #D7D7D7;
    border-radius: 5px;
    background-color: white;
    font-size: 14px;
}
.infoheight
{
    min-height: 10rem;
}
.tag{
    height:20px;
    line-height:normal;
}
</style>
