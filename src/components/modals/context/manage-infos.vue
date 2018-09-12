<template lang="pug">
  v-card
    v-card-text
      .title.text-xs-center.my-2 {{ $t('modals.createEntity.infosList') }}
      v-list(v-if="Object.keys(infos).length")
        v-list-group.my-1(
        v-for="infoName in infos",
        :key="infoName.name"
        )
          v-list-tile(slot="activator")
            v-list-tile-content
              v-list-tile-title {{ infoName.name }}
            v-list-tile-action
              v-btn(icon, flat, @click.stop="deleteInfo(infoName)")
                v-icon delete
          v-list-tile(@click="")
            v-list-tile-content Description : {{ infoName.description }}
            v-list-tile-content Value : {{ infoName.value }}
      v-card-text(v-else) No infos
      v-form(ref="infoForm")
        .title.text-xs-center.my-2 {{ $t('modals.createEntity.addInfos') }}
        v-layout
          v-text-field(
          :label="$t('common.name')",
          v-model="form.name",
          v-validate="'required|unique-name'",
          data-vv-name="name",
          :error-messages="errors.collect('name')"
          )
          v-text-field(
          :label="$t('common.description')",
          v-model="form.description",
          v-validate="'required'",
          data-vv-name="description",
          :error-messages="errors.collect('description')"
          )
          v-text-field(
          :label="$t('common.value')",
          v-model="form.value",
          v-validate="'required'",
          data-vv-name="value",
          :error-messages="errors.collect('value')"
          )
          v-btn(icon, flat, @click="addInfo")
            v-icon add
</template>

<script>
import modalInnerMixin from '@/mixins/modal/modal-inner';
import { MODALS } from '@/constants';

export default {
  name: MODALS.contextInfos,
  $_veeValidate: {
    validator: 'new',
  },
  mixins: [
    modalInnerMixin,
  ],
  props: {
    template: {
      type: String,
    },
  },
  data() {
    return {
      form: {
        name: '',
        description: '',
        value: '',
      },
      infos: {},
    };
  },
  mounted() {
    this.createUniqueValidationRule();
    this.infos = this.config.item ? this.config.item.infos : {};
  },
  methods: {
    async addInfo() {
      const isFormValid = await this.$validator.validateAll();
      if (isFormValid) {
        this.infos[this.form.name] = { ...this.form };
        this.$refs.infoForm.reset();
        this.$validator.reset();
        this.$emit('update:infos', this.infos);
      }
    },
    deleteInfo(name) {
      delete this.infos[name];
      this.$emit('update:infos', this.infos);
      this.$validator.reset();
      this.$forceUpdate();
      this.$nextTick(() => {
        if (this.form.name) {
          this.$validator.validate();
        }
      });
    },
    createUniqueValidationRule() {
      this.$validator.extend('unique-name', {
        getMessage: () => this.$t('validator.unique'),
        validate: value => !this.forbiddenNames().includes(value),
      });
    },
    forbiddenNames() {
      return Object.keys(this.infos);
    },
  },
};
</script>
