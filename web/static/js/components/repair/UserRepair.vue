<template> <div class="grid-content" v-loading="this.$store.state.loading"> <el-row> <el-col :span=12> <el-input v-on:input="OnInput" v-model.number="phone" class="input" clearable> <template slot="prepend">用户ID</template>
                    </el-input>
                <el-col>
            </el-row>
            <el-row class="input">
                <el-button class="button" type="primary" @click="Query()" plain>查询</el-button>
                <el-button class="button" type="primary" @click="Submit()" plain :disabled="notqueryed">删除</el-button>
            </el-row>
            <el-row>
                <pre>{{userDetail}}</pre>
            </el-row>
        </el-col>
    </div>
</template>

<script>
    export default {
        name: "v-userrepair",
        data() {
            return {
                phone : '',
                userDetail : {}
            }
        },
        methods: {
            OnInput(val) {
                this.userDetail = {}
            },
            Query() {
                this.$http
                    .post('http://cgi.gstyun.cn/cgi-bin/user/userquerydetail', {"mobile_number":this.phone})
                    .then(response => {
                        this.userDetail = response.data
                    })
            },
            Submit() {
                this.$store.commit("DeleteUser", {"phone":this.phone})
            }
        },
        computed: {
            notqueryed : {
                get() {
                    if (this.userDetail["status"] == "0") {
                        console.log("false")
                        return false
                    }
                    console.log("true")
                    return true
                }
            }
        }
    }
</script>

<style scoped>
.input {
    border-radius: 4px;
    padding: 10px 0;
    background-color: #f9fafc;
    margin-top: 10px;
    margin-bottom: 10px;
}
pre {}
.string { color: green; }
.number { color: darkorange; }
.boolean { color: blue; }
.null { color: magenta; }
.key { color: red; }
</style>

