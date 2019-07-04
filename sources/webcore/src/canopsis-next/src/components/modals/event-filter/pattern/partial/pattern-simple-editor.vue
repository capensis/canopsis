<template lang="pug">
  div
    v-layout(justify-end)
      v-tooltip(v-for="(action, index) in mainActions", :key="`action-${index}`", top)
        v-btn(slot="activator", icon, @click="action.action()")
          v-icon(:class="action.iconClass") {{ action.icon }}
        span {{ action.tooltip }}
    v-layout(row)
      v-flex(xs12)
        v-treeview(:items="treeViewItems", :open.sync="opened", open-all)
          template(slot="label", slot-scope="{ item }")
            v-flex(xs12)
              v-layout(row)
                v-flex(xs6) {{ item.name }}
                  span(v-show="item.isValueRule") :
                template(v-if="item.isValueRule")
                  v-flex(v-if="isSimpleValueRule(item.value)")
                    span.body-1.font-italic {{ item.value | treeViewValue }}
                  v-flex(v-else)
                    v-layout(column)
                      v-flex(v-for="(field, fieldKey) in item.value", :key="fieldKey")
                        p.body-1.font-italic {{ fieldKey }} {{ field }}
          template(slot="append", slot-scope="{ item }")
            div
              v-tooltip(v-for="(action, index) in getActionsForItem(item)", :key="`action-${index}`", top)
                v-btn(slot="activator", icon, @click="action.action(item)")
                  v-icon(:class="action.iconClass") {{ action.icon }}
                span {{ action.tooltip }}
</template>

<script>
import { isObject, isString, isNull, dropRight, has } from 'lodash';

import { MODALS } from '@/constants';

import formMixin from '@/mixins/form';
import modalMixin from '@/mixins/modal';

export default {
  filters: {
    treeViewValue(value) {
      if (isString(value)) {
        return `"${value}"`;
      } else if (isNull(value)) {
        return 'null';
      }

      return value;
    },
  },
  mixins: [formMixin, modalMixin],
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
          tooltip: this.$t('modals.eventFilterRule.tooltips.addValueRuleField'),
          icon: 'add',
          iconClass: 'primary--text',
          action: this.showAddValueRuleFieldModal,
        },
        editValueRuleField: {
          tooltip: this.$t('modals.eventFilterRule.tooltips.editValueRuleField'),
          icon: 'edit',
          action: this.showEditValueRuleFieldModal,
        },
        addObjectRuleField: {
          tooltip: this.$t('modals.eventFilterRule.tooltips.addObjectRuleField'),
          icon: 'library_add',
          iconClass: 'primary--text',
          action: this.showAddObjectRuleFieldModal,
        },
        editObjectRuleField: {
          tooltip: this.$t('modals.eventFilterRule.tooltips.editObjectRuleField'),
          icon: 'edit',
          action: this.showEditObjectRuleFieldModal,
        },
        removeRuleField: {
          tooltip: this.$t('modals.eventFilterRule.tooltips.removeRuleField'),
          icon: 'remove',
          iconClass: 'error--text',
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

    getActionsForItem() {
      const { actionsMap } = this;

      return (treeViewItem) => {
        if (has(treeViewItem, 'value')) {
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
      };
    },

    treeViewItems() {
      return this.parsePatternToTreeview(this.pattern);
    },

    isSimpleValueRule() {
      return rule => !isObject(rule);
    },

    isValueRule() {
      return (rule) => {
        if (isObject(rule)) {
          const items = Object.entries(rule);

          return items.length && items.every(([key, value]) => this.operators.indexOf(key) !== -1 && !isObject(value));
        }

        return true;
      };
    },
  },
  methods: {
    /**
     * Parse pattern object to treeview items
     *
     * @param {Object} source
     * @param {Array} prevPath
     */
    parsePatternToTreeview(source, prevPath = []) {
      return Object.entries(source).map(([key, value]) => {
        const path = [...prevPath, key];
        const item = {
          path,
          name: key,
          id: path.join('.'),
          isValueRule: this.isValueRule(value),
        };

        if (item.isValueRule) {
          item.value = value;
        } else {
          item.children = this.parsePatternToTreeview(value, path);
        }

        return item;
      }, []);
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

      this.showModal({
        name: MODALS.addEventFilterRuleToPattern,
        config: {
          operators: this.operators,
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
      const { name, value, path } = treeViewItem;

      this.showModal({
        name: MODALS.addEventFilterRuleToPattern,
        config: {
          ruleKey: name,
          ruleValue: value,
          operators: this.operators,
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

      this.showModal({
        name: MODALS.textFieldEditor,
        config: {
          title: this.$t('modals.eventFilterRule.tooltips.addObjectRuleField'),
          field: {
            label: this.$t('modals.eventFilterRule.field'),
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
      this.showModal({
        name: MODALS.textFieldEditor,
        config: {
          title: this.$t('modals.eventFilterRule.tooltips.editObjectRuleField'),
          field: {
            value: treeViewItem.name,
            label: this.$t('modals.eventFilterRule.field'),
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
