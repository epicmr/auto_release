<template>
    <div v-loading="this.$store.state.loading">
        <el-row>
            <el-input class="input" v-on:input="OnInput" v-model.number="phone" placeholder="请输入用户手机号码" @keyup.enter.native="Query()" clearable></el-input>
            <el-button class="button" type="primary" @click="Query()" accesskey="1">查询</el-button>
        </el-row>
        <el-row>
            <el-radio-group class="radiogroup" v-model.number="this.vipLv" v-if="queryed">
                <el-radio :label="0" type="number">无</el-radio>
                <el-radio :label="1" type="number">VIP</el-radio>
                <el-radio :label="2" type="number">SVIP</el-radio>
                <el-radio :label="3" type="number">SSVIP</el-radio>
            </el-radio-group>
        </el-row>
        <div class="infowrapper" v-if="queryed">
            <el-row class="info">
                <el-col :span=4>
                    <span>优惠券码</span>
                </el-col>
                <el-col :span=4>
                    <span>使用状态</span>
                </el-col>
                <el-col :span=12>
                    <span>名称</span>
                </el-col>
            </el-row>
            <el-row class="info infoheight">
                <el-row v-for="r in this.vipRights">
                    <el-col :span=4>
                        <span>{{r.record_id}}</span>
                    </el-col>
                    <el-col :span=4>
                        <el-tag :type="stateTypes[r.use_state]" class="tag" effect="dark">{{stateDescs[r.use_state]}}</el-tag>
                    </el-col>
                    <el-col :span=12>
                        <span>{{r.group_name}}</span>
                    </el-col>
                </el-row>
            </el-row>
        </div>
        <el-row>
            <el-button class="button" type="primary" @click="Submit()" v-if="isVip">解除会员</el-button>
        </el-row>
    </div>
</template>

<script>
    export default {
        name: "v-userrepair",
        data() {
            return {
                userId : 0,
                vipLv : 0,
                vipRights : '',
                phone : '',
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
                this.vipLv = 0
                this.vipRights = ''
                this.userId = 0
            },
            async Query() {
                await this.$http
                    .post('http://cgi.gstyun.cn/cgi-bin/user/userquerydetail', {"mobile_number":this.phone})
                    .then(response => {
                        if (parseInt(response.data.status, 10) == 0) {
                            this.vipLv = response.data.member_level
                            this.userId = response.data.user_id
                        }
                    })
                let payload = {
                    "page_size":"99",
                    "page_no":"1",
                    "user_id":this.userId,
                    "coupon_user_grade": 1 << this.vipLv
                }
                await this.$http
                    .get('http://stg-cgi.gstyun.cn/cgi-bin/coupon/queryusercouponlist', {params:payload})
                    .then(response => {
                        let couponList = response.data.data.coupon_list
                        if (Array.isArray(couponList)) {
                            this.vipRights = couponList
                        }
                    })
            },
            Submit() {
                let payload = {
                    "user_id" : this.userId
                }
                this.$http
                    .post('http://stg-cgi.gstyun.cn/cgi-bin/coupon/disablemembercoupon', payload)
                    .then(response => {
                        if (parseInt(response.data.status, 10) == 0) {
                            this.$message('失效会员优惠券成功');
                        }
                    })
                payload = {
                    "user_id":this.userId,
                    "member_type": "0",
                    "name": "",
                    "discount": "0",
                    "city_id": "0",
                    "type_id":"0"
                }
                this.$http
                    .post('http://stg-cgi.gstyun.cn/cgi-bin/member/updateaccount', payload)
                    .then(response => {
                        if (parseInt(response.data.status, 10) == 0) {
                            this.$message('失效会员身份成功');
                        }
                    })

                setTimeout(this.Query, 200)
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
            },
            isVip : {
                get() {
                    if (this.vipLv > 0) {
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
