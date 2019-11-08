<template>
    <div v-loading="this.$store.state.loading">
        <el-row>
            <el-input class="input" v-on:input="OnInput" v-model.number="phone" placeholder="请输入用户手机号码" @keyup.enter.native="Query()" clearable></el-input>
            <el-button class="button" type="primary" @click="Query()" accesskey="1">查询</el-button>
        </el-row>
        <div class="infowrapper" v-if="queryed">
            <el-row class="info">
                <el-col :span=4>
                    <span>用户ID</span>
                </el-col>
                <el-col :span=4>
                    <span>昵称</span>
                </el-col>
                <el-col :span=4>
                    <span>手机</span>
                </el-col>
                <el-col :span=4>
                    <span>VIP</span>
                </el-col>
            </el-row>
            <el-row class="info infoheight">
                <el-col class="hold" :span=4>
                    <span>{{userDetail.user_id}}</span>
                </el-col>
                <el-col class="hold" :span=4>
                    <span>{{userDetail.nick_name}}</span>
                </el-col>
                <el-col class="hold" :span=4>
                    <span>{{userDetail.mobile_phone}}</span>
                </el-col>
                <el-col class="hold" :span=4>
                    <span>{{userDetail.member_level}}</span>
                </el-col>
            </el-row>
        </div>
        <el-row>
            <el-button class="button" type="danger" @click="Submit()" v-if="queryed">删除</el-button>
        </el-row>
    </div>
</template>

<script>
    export default {
        name: "v-userrepair",
        data() {
            return {
                phone : '',
                userDetail : ''
            }
        },
        methods: {
            OnInput(val) {
                this.userDetail = ''
            },
            Query() {
                this.$http
                    .post('http://cgi.gstyun.cn/cgi-bin/user/userquerydetail', {"mobile_number":this.phone})
                    .then(response => {
                        this.userDetail = response.data
                    })
                    .catch(err => {
                        console.info("%c [axiso catch error]", "color:orange", err)
                    });
            },
            Submit() {
                this.$store.commit("DeleteUser", {"phone":this.phone})
            }
        },
        computed: {
            queryed : {
                get() {
                    if (this.userDetail["status"] == "0") {
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
.hold
{
    height: 15px;
}
</style>
