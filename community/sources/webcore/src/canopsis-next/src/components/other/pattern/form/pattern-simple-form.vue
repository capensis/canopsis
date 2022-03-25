<template lang="pug">
  div.pattern-simple-editor
    v-layout(justify-end)
      c-action-btn(
        v-for="(action, index) in mainActions",
        :key="index",
        :tooltip="action.tooltip",
        :icon="action.icon",
        :color="action.iconColor",
        @click="action.action()"
      )
    v-layout(row)
      pattern-information(v-if="treeViewItems.length > 1") {{ $t('common.and') }}
      v-flex(xs12)
        v-treeview(:items="treeViewItems", :open.sync="opened", open-all)
          template(slot="label", slot-scope="{ item }")
            v-flex(xs12)
              v-layout(row)
                v-flex.text-field(xs6) {{ item.name }}
                  span(v-show="item.rule") :
                template(v-if="item.rule")
                  v-flex(v-if="item.isSimpleRule", xs6)
                    span.body-1.font-italic.text-field {{ item.rule.value | treeViewValue }}
                  v-flex(v-else, xs6)
                    v-flex(v-for="(field, fieldKey) in item.rule.value", :key="fieldKey")
                      p.body-1.font-italic {{ fieldKey }}
                      p.body-1.font-italic.text-field {{ field | treeViewValue }}
          template(slot="append", slot-scope="{ item }")
            v-layout(row)
              c-action-btn(
                v-for="(action, index) in getActionsForItem(item)",
                :key="index",
                :tooltip="action.tooltip",
                :icon="action.icon",
                :color="action.iconColor",
                @click="action.action(item)"
              )
</template>

<script>
import { isString, isNull, dropRight, has } from 'lodash';

import { MODALS } from '@/constants';

import { convertPatternToTreeview } from '@/helpers/treeview';

import { formMixin } from '@/mixins/form';

import PatternInformation from '@/components/other/pattern/pattern-information.vue';

export default {
  components: { PatternInformation },
  filters: {
    treeViewValue(value) {
      if (isString(value)) {
        return `"${value}"`;
      }

      if (isNull(value)) {
        return 'null';
      }

      return value;
    },
  },
  mixins: [formMixin],
  model: {
    prop: 'pattern',
    event: 'input',
  },
  props: {
    pattern: {
      type: Object,
      required: true,
    },
    operators: {
      type: Array,
      required: true,
    },
    onlySimpleRule: {
      type: Boolean,
      default: false,
    },
  },
  data() {
    return {
      opened: [],
    };
  },
  computed: {
    actionsMap() {
      return {
        addValueRuleField: {
          tooltip: this.$t('eventFilter.tooltips.addValueRuleField'),
          icon: 'add',
          iconColor: 'primary',
          action: this.showAddValueRuleFieldModal,
        },
        editValueRuleField: {
          tooltip: this.$t('eventFilter.tooltips.editValueRuleField'),
          icon: 'edit',
          action: this.showEditValueRuleFieldModal,
        },
        addObjectRuleField: {
          tooltip: this.$t('eventFilter.tooltips.addObjectRuleField'),
          icon: 'library_add',
          iconColor: 'primary',
          action: this.showAddObjectRuleFieldModal,
        },
        editObjectRuleField: {
          tooltip: this.$t('eventFilter.tooltips.editObjectRuleField'),
          icon: 'edit',
          action: this.showEditObjectRuleFieldModal,
        },
        removeRuleField: {
          tooltip: this.$t('eventFilter.tooltips.removeRuleField'),
          icon: 'remove',
          iconColor: 'error',
          action: this.deleteRule,
        },
      };
    },

    mainActions() {
      const { actionsMap } = this;

      return [
        actionsMap.addValueRuleField,
        actionsMap.addObjectRuleField,
      ];
    },

    treeViewItems() {
      return convertPatternToTreeview(this.pattern, this.operators);
    },
  },
  methods: {
    getActionsForItem(treeViewItem) {
      const { actionsMap } = this;

      if (has(treeViewItem, 'rule')) {
        return [
          actionsMap.editValueRuleField,
          actionsMap.removeRuleField,
        ];
      }

      return [
        actionsMap.addValueRuleField,
        actionsMap.addObjectRuleField,
        actionsMap.editObjectRuleField,
        actionsMap.removeRuleField,
      ];
    },

    /**
     * Open treeview item
     *
     * @param {Object} treeViewItem
     */
    openTreeviewItem(treeViewItem) {
      if (treeViewItem && this.opened.indexOf(treeViewItem.id) === -1) {
        this.opened.push(treeViewItem.id);
      }
    },

    /**
     * Show modal window for adding of usually rule value field
     *
     * @param {Object} treeViewParent - parent of rule value which we will add
     */
    showAddValueRuleFieldModal(treeViewParent) {
      const parentPath = treeViewParent ? treeViewParent.path : [];

      this.$modals.show({
        name: MODALS.createPatternRule,
        config: {
          operators: this.operators,
          onlySimple: this.onlySimpleRule,
          action: (newRule) => {
            this.updateField([...parentPath, newRule.field], newRule.value);

            this.$nextTick(() => this.openTreeviewItem(treeViewParent));
          },
        },
      });
    },

    /**
     * Show modal window for editing of usually rule value field
     *
     * @param {Object} treeViewItem
     */
    showEditValueRuleFieldModal(treeViewItem) {
      const { rule, path } = treeViewItem;

      this.$modals.show({
        name: MODALS.createPatternRule,
        config: {
          rule,

          operators: this.operators,
          onlySimple: this.onlySimpleRule,
          action: (newRule) => {
            const newPath = [...dropRight(path, 1), newRule.field];

            this.updateAndMoveField(path, newPath, newRule.value);
          },
        },
      });
    },

    /**
     * Show modal window for adding of object wrapper rule field
     *
     * @param {Object} treeViewParent - parent of rule wrapper which we will add
     */
    showAddObjectRuleFieldModal(treeViewParent) {
      const parentPath = treeViewParent ? treeViewParent.path : [];

      this.$modals.show({
        name: MODALS.textFieldEditor,
        config: {
          title: this.$t('eventFilter.tooltips.addObjectRuleField'),
          field: {
            label: this.$t('eventFilter.field'),
            validationRules: 'required',
            name: 'field',
          },
          action: (field) => {
            this.updateField([...parentPath, field], {});

            this.$nextTick(() => this.openTreeviewItem(treeViewParent));
          },
        },
      });
    },

    /**
     * Show modal window for editing of object wrapper rule field
     *
     * @param {Object} treeViewItem
     */
    showEditObjectRuleFieldModal(treeViewItem) {
      this.$modals.show({
        name: MODALS.textFieldEditor,
        config: {
          title: this.$t('eventFilter.tooltips.editObjectRuleField'),
          field: {
            value: treeViewItem.name,
            label: this.$t('eventFilter.field'),
            validationRules: 'required',
            name: 'field',
          },
          action: field => this.moveField([...treeViewItem.path], [field]),
        },
      });
    },

    /**
     * Remove rule field
     *
     * @param {Object} treeViewItem
     */
    deleteRule(treeViewItem) {
      this.removeField(treeViewItem.path);
    },
  },
};
</script>

<style lang="scss" scoped>
  .pattern-simple-editor {
    & /deep/ {
      .v-treeview-node__content, .v-treeview-node__label {
        flex-shrink: 8;
      }
      .v-treeview-node__label {
        width: 100%;
      }
      .text-field {
        word-break: break-all;
        margin-bottom: 0;
      }
    }
  }
</style>
