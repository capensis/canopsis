<template lang="pug">
  div.request-information-content-row(v-click-outside.contextmenu="clickOutsideDirective")
    v-menu(
      v-model="contextmenu.active",
      :position-x="contextmenu.x",
      :position-y="contextmenu.y"
    )
      v-list(dense)
        v-list-tile(v-for="item in items", :key="item.text", @click="item.action")
          v-list-tile-content
            v-list-tile-title {{ item.text }}
    span.request-information-content-row__key(@contextmenu.prevent="openContextmenu") {{ row.name }}
    template(v-if="row.value")
      span :&nbsp;
      span.request-information-content-row__value(@contextmenu.prevent="openContextmenu") {{ row.value }}
</template>

<script>
import { writeTextToClipboard } from '@/helpers/clipboard';

export default {
  props: {
    row: {
      type: Object,
      required: true,
    },
  },
  data() {
    return {
      contextmenu: {
        active: false,
        x: 0,
        y: 0,
      },
    };
  },
  computed: {
    clickOutsideDirective() {
      return {
        handler: this.closeContextMenu,
        closeConditional: () => true,
      };
    },

    items() {
      const items = [];

      if (this.row.value) {
        items.push({
          text: this.$t('common.copyValue'),
          action: this.copyRowValue,
        });
      }

      if (this.row.name) {
        items.push({
          text: this.$t('common.copyProperty'),
          action: this.copyRowProperty,
        });
      }

      if (this.row.path) {
        items.push({
          text: this.$t('common.copyPropertyPath'),
          action: this.copyRowPropertyPath,
        });
      }

      return items;
    },
  },
  methods: {
    openContextmenu(event) {
      this.contextmenu.x = event.x;
      this.contextmenu.y = event.y;
      this.contextmenu.active = true;
    },

    closeContextMenu() {
      this.contextmenu.active = false;
    },

    async copy(value) {
      try {
        await writeTextToClipboard(value);

        this.$popups.success({ text: this.$t('success.valueCopied') });
      } catch (err) {
        this.$popups.error({ text: this.$t('errors.default') });
      }
    },

    copyRowValue() {
      this.copy(this.row.value);
    },

    copyRowProperty() {
      this.copy(this.row.name);
    },

    copyRowPropertyPath() {
      this.copy(this.row.path);
    },
  },
};
</script>

<style lang="scss">
.request-information-content-row {
  &__key, &__value {
    transition: 0.16s;
    border-radius: 4px;
    padding: 2px;
    cursor: pointer;

    &:hover {
      background: rgba(0, 0, 0, 0.1);
    }
  }

  &__key {
    font-weight: bold;
  }
}
</style>
