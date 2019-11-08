<template>
    <div class="grid-content" v-loading="this.loading">
        <el-col>
            <el-row>
                <el-select @change="handleItemChange(selectedName)" v-model:value="selectedName" filterable placeholder="请选择">
                    <el-option v-for="host in this.hosts"
                               :key="host.id"
                               :label="host.name"
                               :value="host.name">
                    </el-option>
                </el-select>
                <el-button class="button" type="primary" @click="EnableEdit" plain>{{AddOrEdit}}</el-button>
            </el-row>

            <el-row>
                <el-col :span=16>
                    <el-checkbox-group class="checkboxgroup" v-model="servTypeList" :disabled="editdisable">
                        <el-checkbox :label="1">CGI</el-checkbox>
                        <el-checkbox :label="2">DAO</el-checkbox>
                        <el-checkbox :label="3">GO</el-checkbox>
                    </el-checkbox-group>
                    <el-select class="input" v-model="selected.env_id" filterable placeholder="请选择环境" :disabled="editdisable">
                        <el-option v-for="env in this.envs" :key="env.id" :label="env.name" :value="env.id"> </el-option>
                    </el-select>
                    <el-input v-model="selected.name" class="input" :disabled="editdisable">
                        <template slot="prepend">主机名称</template>
                    </el-input>
                </el-col>
            </el-row>
            <el-row>
                <el-button class="button" type="primary" @click="Submit()" :disabled="editdisable" plain>提交</el-button>
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
        },
        methods: {
            handleItemChange(val) {
                this.DisableEdit()
                console.log(this.hosts)
                for (let i in this.hosts) {
                    let host = this.hosts[i]
                    if (host.name == this.selectedName) {
                        this.selected = host
                        break
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
                this.DisableEdit()
                this.$store.commit("Host", this.selected)
            }
        },
        computed: {
            ...mapState([
                "loading",
                "envs",
                "hosts"
            ]),
            servTypeList: {
                get() {
                    let l = []
                    for (let i = 0; i < 32; i++) {
                        if (this.selected.serv_type & (1 << i)) {
                            l.push(i)
                        }
                    }
                    return l
                },
                set(val) {
                    this.selected.serv_type = 0
                    for (let i in val) {
                        this.selected.serv_type += (1 << val[i])
                    }
                }
            },
            AddOrEdit: {
                get() {
                    if (this.selected.id == 0 ) {
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
.checkboxgroup {
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

