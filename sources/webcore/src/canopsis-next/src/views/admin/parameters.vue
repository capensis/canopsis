<template lang="pug">
  v-container
    h2.text-xs-center.my-3.display-1.font-weight-medium {{ $t('common.parameters') }}
    v-card.ma-2
      v-card-text
        v-list
          v-list-tile
            v-list-tile-title {{ $t('parameters.groupsNavigationType.title') }}
            v-list-tile-content
              v-select.mt-0(
              :items="groupsNavigationOptions",
              :value="groupsNavigationType",
              @change="setGroupsNavigationType"
              )
    v-layout(row, wrap)
      v-flex(xs6, offset-xs3)
        v-card.ma-2
          v-card-text
            user-interface-form
</template>

<script>
import { GROUPS_NAVIGATION_TYPES } from '@/constants';

import { getFileDataUrlContent } from '@/helpers/file-select';

import appMixin from '@/mixins/app';
import entitiesInfoMixin from '@/mixins/entities/info';

import UserInterfaceForm from '@/components/forms/user-interface.vue';

const MAX_LOGO_SIZE = 16777216;

export default {
  components: { UserInterfaceForm },
  mixins: [appMixin, entitiesInfoMixin],
  data() {
    return {
      userInterfaceForm: {
        appTitle: 'Canopsis',
        footer: '',
        logo: '',
      },
      selectedLogoFileName: '',
      pendingLogo: false,
    };
  },
  computed: {
    groupsNavigationOptions() {
      return [
        {
          text: this.$t('parameters.groupsNavigationType.items.sideBar'),
          value: GROUPS_NAVIGATION_TYPES.sideBar,
        },
        {
          text: this.$t('parameters.groupsNavigationType.items.topBar'),
          value: GROUPS_NAVIGATION_TYPES.topBar,
        },
      ];
    },
    btnProps() {
      return {
        loading: this.pendingLogo,
      };
    },
  },
  mounted() {
    this.userInterfaceForm = {
      appTitle: this.appTitle || 'Canopsis',
      footer: this.footer,
      logo: this.logo,
    };
  },
  methods: {
    async changeLogo(e) {
      const { files } = e.target;
      const [file] = files;

      if (file) {
        if (file.size <= MAX_LOGO_SIZE) {
          this.pendingLogo = true;
          this.userInterfaceForm.logo = await getFileDataUrlContent(file);
          this.selectedLogoFileName = file.name;
          this.pendingLogo = false;
        }
      }
    },
    submit() {

    },
  },
};
</script>
