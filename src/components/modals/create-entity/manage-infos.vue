<template lang="pug">
v-card-text
    v-list(v-if="infos.length")
        v-list-group(
          v-for="info in infos"
          :key="info.name"
        )
          v-list-tile(slot="activator")
            v-list-tile-content
              v-list-tile-title {{ info.name }}
            v-list-tile-action
              v-btn(icon, flat, @click.stop="deleteInfo(info.name)")
                v-icon delete
          v-list-tile(@click="")
            v-list-tile-content
              v-list-tile-title Description : {{ info.description }}
              v-list-tile-title Value : {{ info.value }}
    v-card-text(v-else) No infos
    v-btn(flat, @click="showForm = !showForm") Add info
    v-layout(v-show="showForm")
      v-text-field(
        :label="$t('common.name')",
        v-model="form.name",
        v-validate="'required'",
      )
      v-text-field(
      :label="$t('common.description')",
      v-model="form.description",
      v-validate="'required'",
      )
      v-text-field(
      :label="$t('common.value')",
      v-model="form.value",
      v-validate="'required'",
      )
      v-btn(icon, flat, @click="addInfo")
        v-icon done
</template>

<script>
import filter from 'lodash/filter';
import ModalInnerMixin from '@/mixins/modal/modal-inner';
import { MODALS } from '@/constants';


export default {
  name: MODALS.contextInfos,
  mixins: [ModalInnerMixin],
  props: {
    template: {
      type: String,
    },
  },
  data() {
    return {
      showForm: false,
      infos: [],
      form: {
        name: '',
        description: '',
        value: '',
      },
    };
  },
  mounted() {
    if (this.config) {
      this.form.infos = this.config.item.infos;
      this.$emit('update:infos', this.infos);
    }
  },
  methods: {
    addInfo() {
      this.infos.push({ ...this.form });
      this.$emit('update:infos', this.infos);
    },
    deleteInfo(name) {
      this.infos = filter(this.infos, info => info.name !== name);
      this.$emit('update:infos', this.infos);
    },
  },
};
</script>
