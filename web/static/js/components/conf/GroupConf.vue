<template>
    <div class="grid-content" v-loading="this.loading">
        <el-col>
            <el-row>
                <el-col :span=16>
                    <el-radio-group class="radiogroup" v-model="selected.type" :disabled="editdisable">
                        <el-radio label="g">用户组</el-radio>
                        <el-radio label="g2">资源组</el-radio>
                    </el-radio-group>
                    <el-input v-model="selected.name" class="input" :disabled="editdisable">
                        <template slot="prepend">名称</template>
                    </el-input>
                </el-col>
            </el-row>
            <el-row>
                <el-button class="button" type="primary" @click="EnableEdit" plain>{{AddOrEdit}}</el-button>
                <el-button class="button" type="primary" @click="Submit()" :disabled="editdisable" plain>提交</el-button>
            </el-row>
        </el-col>
    </div>
</template>

<script>
    import { mapState } from 'vuex'

    export default {
        name: "v-groupconf",
        data() {
            return {
                selectedName: '',
                selected: {"type":"g"},
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
            EnableEdit() {
                console.log(this.selected)
                this.editdisable = false
            },
            DisableEdit() {
                console.log(this.selected)
                this.editdisable = true
            },
            Submit() {
                this.DisableEdit()
                let list = []
                list.push(this.selected)
                this.$store.commit("UserGroup", list)
            }
        },
        computed: {
            ...mapState([
                "loading"
            ]),
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
</style>

