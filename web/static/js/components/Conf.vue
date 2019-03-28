<template>
    <div class="grid-content">
        <v-select :list="this.$store.state.confs">
        </v-select>
    </div>
</template>

<script>

import VSelect from "./Select";

export default {
    components: {
        VSelect
    },
    data() {
        return {
            checkedApp:[]
        }
    }
};
</script>
