<template lang="pug">
v-card
  v-card-title.primary.white--text
    v-layout(justify-space-between, align-center)
      h2 {{ config.title }}
  v-card-text.py-0
    v-container
      v-form
        v-text-field(
        v-model="form._id",
        :label="$t('common.name')",
        name="name",
        v-validate="'required'",
        :error-messages="errors.collect('name')"
        )
        v-text-field(v-model="form.description", :label="$t('common.description')")
        v-layout(align-center)
          div Default view:
          div.pl-2 {{ getViewTitle(this.form.defaultview) }}
          v-menu(
          offset-y,
          :close-on-content-click="false",
          v-model="defaultViewMenu"
          )
            v-btn(slot="activator", fab, small, depressed)
              v-icon edit
            v-list.py-0
              v-list-group(v-for="group in groups", :key="group._id")
                v-list-tile(slot="activator")
                  v-list-tile-content
                    v-list-tile-title {{ group.name }}
                v-list-tile(v-for="view in group.views", :key="view._id", @click="updateDefaultView(view)")
                  v-list-tile-content
                    v-list-tile-title {{ view.name }}
    v-divider
    v-layout.py-1(justify-end)
      v-btn(@click="hideModal", depressed, flat) {{ $t('common.cancel') }}
      v-btn.primary.white--text(@click="submit") {{ $t('common.submit') }}
</template>

<script>
import pick from 'lodash/pick';
import { MODALS } from '@/constants';
import { generateRole } from '@/helpers/entities';
import modalInnerMixin from '@/mixins/modal/modal-inner';
import entitiesViewGroupMixin from '@/mixins/entities/view/group';
import entitiesViewMixin from '@/mixins/entities/view';
import entitiesRoleMixin from '@/mixins/entities/role';

export default {
  name: MODALS.createRole,
  $_veeValidate: {
    validator: 'new',
  },
  mixins: [modalInnerMixin, entitiesViewGroupMixin, entitiesViewMixin, entitiesRoleMixin],
  data() {
    const group = this.modal.config.group || { name: '', description: '', defaultView: '' };

    return {
      form: pick(group, ['_id', 'description', 'defaultview']),
      defaultViewMenu: false,
    };
  },
  computed: {
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
      const isFormValid = await this.$validator.validateAll();

      if (isFormValid) {
        const formData = this.isNew ? { ...generateRole() } : { ...this.role };
        formData._id = this.form._id;

        await this.createRole({ data: { ...formData, ...this.form } });
        await this.fetchRolesListWithPreviousParams();

        this.hideModal();
      }
    },
    updateDefaultView(view) {
      this.form.defaultview = view._id;
      this.defaultViewMenu = false;
    },
    getViewTitle(viewId) {
      const view = this.getViewById(viewId);
      if (view) {
        return view.title;
      }
      return null;
    },
  },
};
</script>

