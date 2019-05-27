<template>
    <div class="grid-content" v-loading="this.loading">
        <el-col>
            <el-row>
                <el-select @change="handleItemChange(selectedName)" v-model="selectedName" filterable placeholder="请选择">
                    <el-option v-for="item in confs"
                               :key="item.serv_name"
                               :label="item.serv_name"
                               :value="item.serv_name">
                    </el-option>
                </el-select>
                <el-button class="button" type="primary" @click="EnableEdit" plain>{{AddOrEdit}}</el-button>
            </el-row>

            <el-row>
                <el-col :span=16>
                    <el-radio-group class="radiogroup" v-model.number="selected.serv_type" :disabled="editdisable">
                        <el-radio :label="1" type="number">CGI</el-radio>
                        <el-radio :label="2" type="number">AO</el-radio>
                        <el-radio :label="3" type="number">DAO</el-radio>
                    </el-radio-group>
                    <el-input class="input" v-model="selected.serv_name" :disabled="editdisable">
                        <template slot="prepend">服务名称</template>
                    </el-input>
                </el-col>
            </el-row>
            <el-row>
                <el-col :span=16>
                    <el-input class="input" v-model="selected.local_path" :disabled="editdisable">
                        <template slot="prepend">本地目录</template>
                    </el-input>
                </el-col>
                <el-col :span=8>
                    <small>{{selected.serv_md5}}</small>
                </el-col>
            </el-row>

            <el-row v-for="env in selected.serv_envs">
                <el-col :span=16>
                    <el-input v-model:value="env.remote_path" class="input" :disabled="editdisable">
                        <template slot="prepend">{{env.env}}</template>
                    </el-input>
                </el-col>
                <el-col :span=8>
                    <small>{{env.serv_md5}}</small>
                </el-col>
            </el-row>

            <el-row>
                <el-button class="button" type="primary" @click="Submit()" :disabled="editdisable" plain>提交</el-button>
                <el-button class="button" type="primary" @click="Refresh()" plain>刷新</el-button>
            </el-row>
        </el-col>
    </div>
</template>

<script>
    import { mapState } from 'vuex'

    export default {
        name: "v-select",
        data() {
            return {
                selectedName: '',
                selected: {"id":0, "serv_type":0},
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
                this.DisableEdit()
                for (let i in this.confs) {
                    let serv = this.confs[i]
                    if (serv.serv_name == this.selectedName) {
                        this.selected = serv
                        break;
                    }
                }
            },
            EnableEdit() {
                this.editdisable = false
            },
            DisableEdit() {
                this.editdisable = true
            },
            Submit() {
                console.log(this.selected)
                this.DisableEdit()
                this.$store.commit("Conf", this.selected)
            },
            Refresh() {
                this.DisableEdit()
                this.$http
                    .get('/api/refresh', {params:{"serv_name": this.selected.serv_name}})
                    .then(response => {
                        let serv = response.data.data
                        if (serv.serv_name == this.selectedName ) {
                            this.selected = serv
                        }
                    })
            }
        },
        computed: {
            ...mapState([
                "loading",
                "confs"
            ]),
            AddOrEdit: {
                get() {
                    if (this.selected.id == 0 ) {
                        return "新增"
                    }
                    return "编辑"
                }
            },
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
.radiogroup {
    border-radius: 4px;
    padding: 10px 0;
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

