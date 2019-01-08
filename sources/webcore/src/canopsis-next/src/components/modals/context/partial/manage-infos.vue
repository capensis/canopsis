<template lang="pug">
  v-card
    v-card-text
      .title.text-xs-center.my-2 {{ $t('modals.createEntity.infosList') }}
      v-list(v-if="infosNames.length")
        v-list-group.my-0(
        v-for="(info, infoName) in infos",
        :key="infoName"
        )
          v-list-tile.py-0(slot="activator")
            v-list-tile-content
              v-list-tile-title {{ infoName }}
            v-list-tile-action
              v-layout
                v-btn.mx-1.primary--text(icon, @click.stop="editInfo(info)")
                  v-icon edit
                v-btn.mx-1.error--text(icon, @click.stop="removeField(infoName)")
                  v-icon delete
          v-list-tile(@click="")
            v-list-tile-content
              v-list-tile-title {{ $t('common.description') }} : {{ info.description }}
              v-list-tile-title {{ $t('common.value') }} : {{ info.value }}
        v-divider
      v-card-text.text-xs-center(v-else) {{ $t('modals.createEntity.noInfos') }}
      v-form(
        ref="form"
      )
        .title.text-xs-center.my-2 {{ $t('modals.createEntity.addInfos') }}
        v-text-field(
        :label="$t('common.name')",
        v-model="form.name",
        v-validate="'required'",
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
        v-textarea(
        :label="$t('common.value')",
        v-model="form.value",
        v-validate,
        data-vv-rule="'required'",
        data-vv-name="value",
        :error-messages="errors.collect('value')"
        )
        v-alert(:value="uniqueInfoError", error) Nom unique
        v-btn(@click="submit", depressed) {{ $t('common.add') }}
</template>

<script>
import find from 'lodash/find';
import formMixin from '@/mixins/form';

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
      form: {
        name: '',
        description: '',
        value: '',
      },
      isEditing: false,
      editingInfoName: '',
      uniqueInfoError: false,
    };
  },
  computed: {
    infosNames() {
      return Object.keys(this.infos);
    },
  },
  methods: {
    async addInfo() {
      const isFormValid = await this.$validator.validateAll();
      const existingInfo = find(this.infos, info => info.name === this.form.name);
      if (isFormValid) {
        if (existingInfo) {
          this.uniqueInfoError = true;
        } else if (this.isEditing) {
          await this.removeField(this.editingInfoName);
          await this.updateField(this.form.name, { ...this.form });
          this.isEditing = false;
          this.uniqueInfoError = false;
          this.$refs.form.reset();
        } else {
          await this.updateField(this.form.name, { ...this.form });
          this.uniqueInfoError = false;
          this.$refs.form.reset();
        }
      }
    },

    async submit() {
      await this.addInfo();
      this.$validator.errors.clear();
    },

    editInfo(info) {
      this.isEditing = true;
      this.editingInfoName = info.name;
      this.form = { ...info };
    },
  },
};
</script>
