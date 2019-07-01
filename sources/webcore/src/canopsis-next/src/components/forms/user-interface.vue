<template lang="pug">
  v-form(@submit="submit")
    v-list
      v-list-tile
        v-list-tile-title App title
        v-list-tile-content
          v-text-field(
          v-model="userInterfaceForm.appTitle",
          label="App title"
          )
      v-list-tile
        v-list-tile-title Login page footer
        v-list-tile-content
          v-text-field(
          v-model="userInterfaceForm.footer",
          label="Footer"
          )
      v-list-tile
        v-list-tile-title Logo
        v-list-tile-content
          file-select(
          ref="fileSelect",
          :btnProps="btnProps",
          tooltip="Select logo file",
          @change="changeLogo"
          )
          v-text-field(
          v-model="userInterfaceForm.logo",
          label="Logo"
          )
    v-alert(
    :value="errors.has('logo')"
    )
    v-divider
    v-btn.mt-3.primary(type="submit") {{ $t('common.submit') }}
</template>

<script>
import { MAX_LOGO_SIZE } from '@/constants';

import uid from '@/helpers/uid';
import { getFileDataUrlContent } from '@/helpers/file-select';

import entitiesInfoMixin from '@/mixins/entities/info';

export default {
  $_veeValidate: {
    validator: 'new',
  },
  mixins: [entitiesInfoMixin],
  data() {
    return {
      logoErrorId: uid(),
      pendingLogoConvert: false,
      selectedLogoFileName: null,
      form: {
        appTitle: 'Canopsis',
        footer: '',
        logo: null,
      },
    };
  },
  methods: {
    async changeLogo(e) {
      const { files } = e.target;
      const [file] = files;

      if (file) {
        if (file.size <= MAX_LOGO_SIZE) {
          this.pendingLogoConvert = true;
          this.userInterfaceForm.logo = await getFileDataUrlContent(file);
          this.selectedLogoFileName = file.name;
          this.pendingLogoConvert = false;
        } else {
          this.errors.add({
            field: 'logo',
            id: this.logoErrorId,
            msg: 'The file size should be less then 16MB',
          });
        }
      }
    },

    submit() {

    },
  },
};
</script>
