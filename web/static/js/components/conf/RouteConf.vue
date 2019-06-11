<template>
    <div class="grid-content" v-loading="this.loading">
        <el-col>
            <el-row>
                <el-col :span=4>
                    <el-menu name="route" @select="SelectHandler" class="el-menu-demo" :unique-opened=true>
                        <v-menu v-bind:items="this.itemstree"></v-menu>
                    </el-menu>
                </el-col>
                <el-col :span=12>
                    <el-row>
                        <el-col :span=8>
                            <el-input v-model.number="selected.parent_id" type="number" class="input" :disabled="editdisable">
                                <template slot="prepend">父节点</template>
                            </el-input>
                        </el-col>
                        <el-col :span=12>
                            <el-select @change="handleItemChange(selectedName)" v-model:value="selectedName" filterable class="input" :disabled="editdisable">
                                <el-option v-for="item in this.items"
                                           :key="item.id"
                                           :label="item.name + item.index"
                                           :value="item.id">
                                </el-option>
                            </el-select>
                        </el-col>
                    </el-row>
                    <el-input v-model="selected.name" class="input" :disabled="editdisable">
                <template slot="prepend">名称</template>
                    </el-input>
                    <el-input v-model="selected.index" class="input" :disabled="editdisable">
                        <template slot="prepend">index</template>
                    </el-input>
                    <el-input v-model="selected.class" class="input" :disabled="editdisable">
                        <template slot="prepend">class</template>
                    </el-input>
                    <el-button class="button" type="primary" @click="EnableEdit" plain>{{AddOrEdit}}</el-button>
                    <el-button class="button" type="primary" @click="Submit()" :disabled="editdisable" plain>提交</el-button>
                </el-col>
            </el-row>
        </el-col>
    </div>
</template>

<script>
import { mapState } from 'vuex'
import VMenu from "../Menu";

    export default {
        name: "v-select",
        components: {
            VMenu
        },
        data() {
            return {
                selectedName: '',
                selected: {"id":0},
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
                console.log(val)
                this.selected.parent_id = val
                console.log(this.selected)
            },
            SelectHandler(index, indexpath) {
                console.log(index,indexpath)
                let m = {}
                for (let i in this.items) {
                    let item = this.items[i]
                    m[item.id] = item.name
                    console.log(item.index, index)
                    if (item.index == index) {
                        console.log(item)
                        this.selected = item
                    }
                }
                this.selectedName = m[this.selected.parent_id]
                console.log(this.selected)
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
                this.$store.commit("Item", this.selected)
            }
        },
        computed: {
            ...mapState([
                "loading",
                "itemstree",
                "items"
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
    padding: 10px 0;
    background-color: #f9fafc;
    margin-top: 10px;
    margin-bottom: 10px;
}
</style>

