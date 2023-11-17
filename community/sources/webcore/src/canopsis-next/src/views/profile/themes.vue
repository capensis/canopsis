<template>
  <c-page
    :creatable="hasCreateAnyThemeAccess"
    :create-tooltip="$t('modals.createTheme.create.title')"
    @refresh="fetchList"
    @create="showCreateThemeModal"
  >
    <themes-list
      :themes="themes"
      :pending="themesPending"
      :options.sync="options"
      :total-items="themesMeta.total_count"
      :updatable="hasUpdateAnyThemeAccess"
      :removable="hasDeleteAnyThemeAccess"
      :duplicable="hasCreateAnyThemeAccess"
      @edit="showEditThemeModal"
      @remove="showRemoveThemeModal"
      @duplicate="showDuplicateThemeModal"
      @remove-selected="showDeleteSelectedThemesModal"
    />
  </c-page>
</template>

<script>
import { omit } from 'lodash';

import { MODALS } from '@/constants';

import { pickIds } from '@/helpers/array';

import { localQueryMixin } from '@/mixins/query-local/query';
import { entitiesThemesMixin } from '@/mixins/entities/theme';
import { permissionsTechnicalProfileThemeMixin } from '@/mixins/permissions/technical/profile/theme';

import ThemesList from '@/components/other/theme/themes-list.vue';

export default {
  inject: ['$system'],
  components: { ThemesList },
  mixins: [
    localQueryMixin,
    entitiesThemesMixin,
    permissionsTechnicalProfileThemeMixin,
  ],
  mounted() {
    this.fetchList();
  },
  methods: {
    showCreateThemeModal() {
      this.$modals.show({
        name: MODALS.createTheme,
        config: {
          action: async (newTheme) => {
            await this.createTheme({ data: newTheme });

            return this.fetchList();
          },
        },
      });
    },

    showEditThemeModal(theme) {
      this.$modals.show({
        name: MODALS.createTheme,
        config: {
          theme,
          title: this.$t('modals.createTheme.edit.title'),
          action: async (newTheme) => {
            await this.updateTheme({ id: theme._id, data: newTheme });

            if (this.currentUser.ui_theme._id === theme._id) {
              this.$system.setTheme(newTheme);
            }

            return this.fetchList();
          },
        },
      });
    },

    showRemoveThemeModal(id) {
      this.$modals.show({
        name: MODALS.confirmation,
        config: {
          action: async () => {
            await this.removeTheme({ id });

            return this.fetchList();
          },
        },
      });
    },

    showDuplicateThemeModal(theme) {
      this.$modals.show({
        name: MODALS.createTheme,
        config: {
          theme: omit(theme, ['_id']),
          title: this.$t('modals.createTheme.duplicate.title'),
          action: async (newTheme) => {
            await this.createTheme({ data: newTheme });

            return this.fetchList();
          },
        },
      });
    },

    showDeleteSelectedThemesModal(selected) {
      this.$modals.show({
        name: MODALS.confirmation,
        config: {
          action: async () => {
            await this.bulkRemoveThemes({ data: pickIds(selected) });

            return this.fetchList();
          },
        },
      });
    },

    fetchList() {
      const params = this.getQuery();

      return this.fetchThemesList({ params });
    },
  },
};
</script>
