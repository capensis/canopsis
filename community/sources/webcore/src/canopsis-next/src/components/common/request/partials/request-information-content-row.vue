<template>
  <div
    v-click-outside.contextmenu="clickOutsideDirective"
    class="request-information-content-row"
  >
    <v-menu
      v-model="contextmenu.active"
      :position-x="contextmenu.x"
      :position-y="contextmenu.y"
    >
      <v-list dense>
        <v-list-item
          v-for="item in items"
          :key="item.text"
          @click="item.action"
        >
          <v-list-item-content>
            <v-list-item-title>{{ item.text }}</v-list-item-title>
          </v-list-item-content>
        </v-list-item>
      </v-list>
    </v-menu>
    <span
      class="request-information-content-row__key"
      @contextmenu.prevent="openContextmenu"
    >
      {{ row.name }}
    </span>
    <template v-if="hasRowValue">
      :&nbsp;
      <span
        class="request-information-content-row__value"
        @contextmenu.prevent="openContextmenu"
      >
        {{ parsedValue }}
      </span>
    </template>
  </div>
</template>

<script>
import { isNull, isString, isUndefined } from 'lodash';

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
    hasRowValue() {
      return !isUndefined(this.row.value);
    },

    parsedValue() {
      if (isNull(this.row.value)) {
        return 'null';
      }

      if (isString(this.row.value) && !this.row.value.length) {
        return '""';
      }

      return this.row.value;
    },

    clickOutsideDirective() {
      return {
        handler: this.closeContextMenu,
        closeConditional: () => true,
      };
    },

    items() {
      const items = [];

      if (this.hasRowValue) {
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
    word-wrap: break-word;

    &:hover {
      background: rgba(0, 0, 0, 0.1);
    }
  }

  &__key {
    font-weight: bold;
  }
}
</style>
