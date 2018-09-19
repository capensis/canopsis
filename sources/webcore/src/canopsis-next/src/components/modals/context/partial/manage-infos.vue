<template lang="pug">
  v-card
    v-card-text
      .title.text-xs-center.my-2 {{ $t('modals.createEntity.infosList') }}
      v-list(v-if="infosNames.length")
        v-list-group.my-1(
        v-for="(info, infoName) in infos",
        :key="infoName"
        )
          v-list-tile(slot="activator")
            v-list-tile-content
              v-list-tile-title {{ infoName }}
            v-list-tile-action
              v-btn(icon, flat, @click.stop="removeField(infoName)")
                v-icon delete
          v-list-tile(@click="")
            v-list-tile-content
              v-list-tile-title {{ $t('common.description') }} : {{ info.description }}
              v-list-tile-title {{ $t('common.value') }} : {{ info.value }}
      v-card-text(v-else) {{ $t('modals.createEntity.noInfos') }}
      v-form
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
import formMixin from '@/mixins/form';

const getDefaultFormData = () => ({
  name: '',
  description: '',
  value: '',
});

/**
 * Form to manipulation with infos
 *
 * @prop {Object} infos - infos from parent
 */
export default {
  $_veeValidate: {
    validator: 'new',
  },
  mixins: [
    formMixin,
  ],
  model: {
    prop: 'infos',
    event: 'input',
  },
  props: {
    infos: {
      type: Object,
      default: () => ({}),
    },
  },
  data() {
    return {
      form: getDefaultFormData(),
    };
  },
  computed: {
    infosNames() {
      return Object.keys(this.infos);
    },
  },
  created() {
    this.createUniqueValidationRule();
  },
  methods: {
    async addInfo() {
      const isFormValid = await this.$validator.validateAll();

      if (isFormValid) {
        this.updateField(this.form.name, { ...this.form });
        this.resetForm();
      }
    },

    createUniqueValidationRule() {
      this.$validator.extend('unique-name', {
        getMessage: () => this.$t('validator.unique'),
        validate: value => !this.infosNames.includes(value),
      });
    },

    resetForm() {
      this.form = getDefaultFormData();
    },
  },
};
</script>
