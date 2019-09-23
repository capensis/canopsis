<template lang="pug">
  v-card(data-test="createRoleModal")
    v-card-title.primary.white--text
      v-layout(justify-space-between, align-center)
        h2 {{ title }}
    v-card-text.py-0
      v-container
        v-form
          v-layout
            v-text-field(
              v-model="form._id",
              :label="$t('common.name')",
              name="name",
              v-validate="'required'",
              :error-messages="errors.collect('name')",
              data-test="name"
            )
          v-layout
            v-text-field(
              v-model="form.description",
              :label="$t('common.description')",
              data-test="description"
            )
          v-layout
            view-selector(v-model="form.defaultview")
      v-divider
      v-layout.py-1(justify-end)
        v-btn(@click="$modals.hide", depressed, flat) {{ $t('common.cancel') }}
        v-btn.primary.white--text(data-test="submitButton", @click="submit") {{ $t('common.submit') }}
</template>

<script>
import { pick } from 'lodash';

import { MODALS } from '@/constants';

import { generateRole } from '@/helpers/entities';

import modalInnerMixin from '@/mixins/modal/inner';
import entitiesViewMixin from '@/mixins/entities/view';
import entitiesRoleMixin from '@/mixins/entities/role';
import entitiesViewGroupMixin from '@/mixins/entities/view/group';

import ViewSelector from './partial/view-selector.vue';

export default {
  name: MODALS.createRole,
  $_veeValidate: {
    validator: 'new',
  },
  components: {
    ViewSelector,
  },
  mixins: [
    modalInnerMixin,
    entitiesViewMixin,
    entitiesRoleMixin,
    entitiesViewGroupMixin,
  ],
  data() {
    const group = this.modal.config.group || { name: '', description: '', defaultView: '' };

    return {
      form: pick(group, ['_id', 'description', 'defaultview']),
      defaultViewMenu: false,
    };
  },
  computed: {
    title() {
      return this.config.title || this.$t('modals.createRole.title');
    },

    role() {
      return this.config.roleId ? this.getRoleById(this.config.roleId) : null;
    },

    isNew() {
      return !this.role;
    },
  },
  mounted() {
    if (!this.isNew) {
      this.form = pick(this.role, [
        '_id',
        'description',
        'defaultview',
      ]);
    }
  },
  methods: {
    async submit() {
      try {
        const isFormValid = await this.$validator.validateAll();

        if (isFormValid) {
          const formData = this.isNew ? { ...generateRole() } : { ...this.role };
          formData._id = this.form._id;

          await this.createRole({ data: { ...formData, ...this.form } });
          await this.fetchRolesListWithPreviousParams();

          this.$popups.addSuccess({ text: this.$t('success.default') });
          this.$modals.hide();
        }
      } catch (err) {
        this.$popups.addError({ text: this.$t('errors.default') });
      }
    },
  },
};
</script>

