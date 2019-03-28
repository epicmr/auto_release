<template>
    <div class="grid-content" v-loading="this.$store.state.loading">
        <el-col>
            <el-row>
                <el-select @change="handleItemChange(value8)" v-model="value8" filterable placeholder="请选择">
                    <el-option v-for="item in list"
                               :key="item.serv.serv_name"
                               :label="item.serv.serv_name"
                               :value="item.serv.serv_name">
                    </el-option>
                </el-select>
                <el-button class="button" type="primary" @click="EnableEdit" plain>{{AddOrEdit}}</el-button>
            </el-row>

            <el-row>
                <el-input class="input" v-model="l_default.serv.serv_type" :disabled="editdisable">
                    <template slot="prepend">服务类型</template>
                </el-input>
                <el-input class="input" v-model="l_default.serv.serv_name" :disabled="editdisable">
                    <template slot="prepend">服务名称</template>
                </el-input>
            </el-row>
            <el-row>
                <el-col :span=16>
                    <el-input class="input" v-model="l_default.serv.local_path" :disabled="editdisable">
                        <template slot="prepend">本地目录</template>
                    </el-input>
                </el-col>
                <el-col :span=8>
                    <small>{{l_default.serv.serv_md5}}</small>
                </el-col>
            </el-row>

            <el-row v-for="(path,env) in l_default.servenv_list">
                <el-col :span=16>
                    <el-input v-model:value="path.remote_path" class="input" :disabled="editdisable">
                        <template slot="prepend">{{env}}</template>
                    </el-input>
                </el-col>
                <el-col :span=8>
                    <small>{{path.serv_md5}}</small>
                </el-col>
            </el-row>

            <el-row>
                <el-button class="button" type="primary" @click="Conf()" :disabled="editdisable" plain>提交</el-button>
                <el-button class="button" type="primary" @click="Refresh()" plain>刷新</el-button>
            </el-row>
        </el-col>
    </div>
</template>

<script>
    export default {
        name: "v-select",
        data() {
            return {
                value8: '',
                l_default: {
                    serv:{},
                    servenv_list:{
                    }
                },
                editdisable: true
            }
        },
        props: {
            list: {
                type: Array,
                default: []
            }
        },
        methods: {
            handleItemChange(val) {
                console.log(this.list)
                console.log(this.value8)
                this.DisableEdit()
                var i
                for (i in this.list) {
                    var l = this.list[i]
                    if (l.serv.serv_name == this.value8 ) {
                        this.l_default = l
                        console.log(this.l_default.serv.serv_name)
                        break;
                    }
                }
            },
            EnableEdit() {
                //this.$store.dispatch("EnableEdit")
                this.editdisable = false
            },
            DisableEdit() {
                //this.$store.dispatch("DisableEdit")
                this.editdisable = true
            },
            Conf() {
                console.log(this.l_default)
                this.DisableEdit()
                this.$store.dispatch("Conf", this.l_default)
                console.log("mm",this.list)
            },
            Refresh() {
                console.log(this.l_default)
                this.DisableEdit()
                this.$http
                    .get('/api/refresh', {params:{"serv_name": this.l_default.serv.serv_name}})
                    .then(response => {
                        var i
                        console.log(response.data)
                        for (i in response.data) {
                            var l = response.data[i]
                            if (l.serv.serv_name == this.value8 ) {
                                this.l_default = l
                                console.log(this.value8)
                                console.log(this.l_default.serv.serv_md5)
                                break;
                            }
                        }
                        console.log(this.l_default)
                    })
            }
        },
        computed: {
            selected: {
                get() {
                },
                set() {
                }
            },
            AddOrEdit: {
                get() {
                    if (this.l_default.serv_name == "" ) {
                        return "新增"
                    }
                    return "编辑"
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
small{
    font-size: 12px;
    color: #66b1ff;
    vertical-align: bottom;
    display: inline-block;
    margin-top: 10px;
}
</style>

