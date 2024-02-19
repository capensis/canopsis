<template>
  <div class="entity-dependencies-by-state-settings">
    <c-zoom-overlay>
      <c-progress-overlay :pending="!ready || pending" />
      <network-graph
        ref="networkGraph"
        :options="options"
        :tooltip-options="tooltipOptions"
        :node-html-label-options="nodeHtmlLabelsOptions"
        :class="{ 'entity-dependencies-by-state-settings-network-graph--ready': ready }"
        class="entity-dependencies-by-state-settings-network-graph fill-height black--text"
        ctrl-wheel-zoom
      />
    </c-zoom-overlay>
  </div>
</template>

<script>
import { omit } from 'lodash';
import { createNamespacedHelpers } from 'vuex';

import { PAGINATION_LIMIT, VUETIFY_ANIMATION_DELAY } from '@/config';
import {
  ROOT_CAUSE_DIAGRAM_OPTIONS,
  ROOT_CAUSE_DIAGRAM_LAYOUT_OPTIONS,
  JUNIT_STATE_SETTING_METHODS,
  ROOT_CAUSE_DIAGRAM_TOOLTIP_OFFSET,
} from '@/constants';

import { normalizeTreeOfDependenciesMapEntities } from '@/helpers/entities/map/list';
import {
  getBadgeElement,
  getButtonHTML,
  getEntityNodeElement,
  getIconElement,
  getProgressElement,
} from '@/helpers/entities/entity/cytoscape';

import { entitiesEntityDependenciesMixin } from '@/mixins/entities/entity-dependencies';

import NetworkGraph from '@/components/common/chart/network-graph.vue';

const { mapActions: mapEntityActions } = createNamespacedHelpers('entity');

export default {
  components: { NetworkGraph },
  mixins: [entitiesEntityDependenciesMixin],
  props: {
    entity: {
      type: Object,
      required: true,
    },
    columns: {
      type: Array,
      default: () => [],
    },
  },
  data() {
    return {
      ready: false,
      pending: true,
      stateSetting: {},
      metaByEntityId: {},
      entitiesById: normalizeTreeOfDependenciesMapEntities([{ entity: this.entity, pinned_entities: [] }]),
    };
  },
  computed: {
    isEventsStateSettings() {
      return !this.stateSetting?.title;
    },

    entitiesElements() {
      const rootElement = this.entitiesById[this.entity._id];
      const { entity, dependencies = [] } = rootElement;

      const elements = [
        {
          group: 'nodes',
          data: {
            id: entity._id,
            entity,
            root: true,
            opened: true,
          },
        },
      ];

      if (this.isEventsStateSettings) {
        elements.push(...this.getEventsNodeElementByEntity(entity));

        return elements;
      }

      elements.push(...this.getEntityDependenciesElement(entity, dependencies, [entity._id]));

      return elements;
    },

    styleOption() {
      return [
        {
          selector: 'node',
          style: {
            width: ROOT_CAUSE_DIAGRAM_OPTIONS.nodeSize,
            height: ROOT_CAUSE_DIAGRAM_OPTIONS.nodeSize,
          },
        },
        {
          selector: 'node[showMore]',
          style: {
            'background-opacity': 0,
            'border-width': 0,
            width: 128,
            height: 34,
          },
        },
        {
          selector: 'node[isEvents]',
          style: {
            width: 30,
            height: 30,
          },
        },
        {
          selector: 'edge',
          style: {
            width: 2,
            'curve-style': 'bezier',
            'line-color': 'silver',
          },
        },
      ];
    },

    nodeHtmlLabelsOptions() {
      return [
        {
          query: 'node',
          valign: 'center',
          halign: 'center',
          tpl: this.getNodeContent,
        },
        {
          query: 'node[showMore]',
          valign: 'center',
          halign: 'center',
          tpl: this.getShowMoreButtonContent,
        },
      ];
    },

    options() {
      const options = {
        ...omit(ROOT_CAUSE_DIAGRAM_OPTIONS, ['nodeSize']),

        style: this.styleOption,
        elements: this.entitiesElements,
      };

      if (this.entitiesElements.length) {
        options.layout = {
          ...ROOT_CAUSE_DIAGRAM_LAYOUT_OPTIONS,
        };
      }

      return options;
    },

    tooltipOptions() {
      return {
        offsetY: (ROOT_CAUSE_DIAGRAM_OPTIONS.nodeSize / 2) + ROOT_CAUSE_DIAGRAM_TOOLTIP_OFFSET,
        getContent: ({ isEvents, entity, root }) => {
          if (isEvents || root) {
            return '';
          }

          return entity.state_setting?.title
            || this.$t(`stateSetting.junit.methods.${JUNIT_STATE_SETTING_METHODS.worst}`);
        },
      };
    },
  },
  watch: {
    entity() {
      this.entitiesById = normalizeTreeOfDependenciesMapEntities([{ entity: this.entity, pinned_entities: [] }]);

      /**
       * TODO: investigate this behavior in the future
       */
      setTimeout(() => this.resetLayout(), 1000);
    },
  },
  async mounted() {
    this.pending = true;
    this.$refs.networkGraph.$cy.on('tap', this.tapHandler);

    await this.fetchEntityStateSetting();

    /**
     * @desc: We are waiting modal showing animation
     */
    setTimeout(() => {
      this.$refs.networkGraph.$cy.center();
      this.ready = true;
    }, VUETIFY_ANIMATION_DELAY);

    if (!this.isEventsStateSettings) {
      await this.fetchDependencies(this.entity._id);
    } else {
      this.runLayout();
    }

    this.pending = false;
  },
  beforeDestroy() {
    this.$refs.networkGraph.$cy.off('tap', this.tapHandler);
  },
  methods: {
    ...mapEntityActions({
      fetchEntityStateSettingWithoutStore: 'fetchStateSettingWithoutStore',
    }),

    async fetchEntityStateSetting() {
      try {
        this.stateSetting = await this.fetchEntityStateSettingWithoutStore({ params: { _id: this.entity._id } });
      } catch (err) {
        console.error(err);
      }
    },

    getEventsNodeElementByEntity(entity) {
      const eventsNodeId = `${entity._id}_events-node`;

      return [
        {
          group: 'nodes',
          data: {
            entity,
            id: eventsNodeId,
            isEvents: true,
          },
        },
        {
          group: 'edges',
          data: {
            source: entity._id,
            target: eventsNodeId,
          },
        },
      ];
    },

    getEntityDependenciesElement(entity, dependenciesIds = [], handledDependenciesIds = []) {
      const dependenciesNodes = dependenciesIds.reduce((acc, childId) => {
        const { dependencies: childDependenciesIds = [], entity: child } = this.entitiesById[childId];

        const isCycle = handledDependenciesIds.includes(childId);

        const hasDependencies = !!childDependenciesIds.length;

        if (!isCycle) {
          const childDependencies = this.getEntityDependenciesElement(
            child,
            childDependenciesIds,
            [...handledDependenciesIds, childId],
          );

          acc.push(
            {
              group: 'nodes',
              data: {
                id: childId,
                entity: child,
                opened: hasDependencies,
              },
            },
            ...childDependencies,
          );
        }

        acc.push(
          {
            group: 'edges',
            data: {
              source: entity._id,
              target: childId,
            },
          },
        );

        if (false) {
          acc.push(...this.getEventsNodeElementByEntity(child));
        }

        return acc;
      }, []);

      dependenciesNodes.push(
        ...this.getShowMoreElements(entity),
      );

      return dependenciesNodes;
    },

    getNodeContent(node) {
      const { entity, pending, opened, root } = node;

      const element = getEntityNodeElement(node);

      /**
       * TODO: Should be changed on state settings deps count
       */
      if (pending || (!root && entity.depends_count > 0)) {
        const badge = getBadgeElement();
        badge.dataset.id = entity._id;

        badge.appendChild(
          pending ? getProgressElement() : getIconElement(opened ? 'remove' : 'add'),
        );

        element.appendChild(badge);
      }

      return element.outerHTML;
    },

    getShowMoreButtonContent(node) {
      const { entity } = node;
      const meta = this.metaByEntityId[entity._id] ?? {};

      const fetchedEntities = meta.page * meta.per_page;

      return getButtonHTML(
        this.$t('common.showMore', { current: fetchedEntities, total: meta.total_count }),
      );
    },

    getShowMoreElements(entity) {
      const meta = this.metaByEntityId[entity._id];

      if (!meta || meta.page >= meta.page_count) {
        return [];
      }

      const showMoreId = `show-all-${entity._id}`;

      return [
        {
          group: 'nodes',
          data: {
            id: showMoreId,
            entity,
            showMore: true,
          },
        },
        {
          group: 'edges',
          data: {
            id: `show-all-edge-${entity._id}`,
            source: entity._id,
            target: showMoreId,
          },
        },
      ];
    },

    /**
     * Remove old elements and add new elements to network graph
     */
    resetLayout() {
      if (!this.$refs?.networkGraph?.$cy) {
        return;
      }

      this.$refs.networkGraph.$cy.elements().remove();
      this.$refs.networkGraph.$cy.add(this.entitiesElements);
      this.runLayout();
    },

    /**
     * Run 'cise' layout for rerender clusters
     */
    async runLayout() {
      if (this.$refs.networkGraph.$cy.nodes().empty()) {
        return;
      }

      try {
        await this.$nextTick();

        this.$refs.networkGraph.$cy.layout({
          ...ROOT_CAUSE_DIAGRAM_LAYOUT_OPTIONS,
          animate: false,
        }).run();
      } catch (err) {
        console.warn(err);
      }
    },

    /**
     * Show dependencies for node
     *
     * @param {Object} target
     * @returns {Promise<void>}
     */
    async showDependencies(target) {
      const { id } = target.data();
      const { page } = this.metaByEntityId[id] ?? {};
      const newPage = page ? page + 1 : 1;

      target.data({
        pending: true,
      });

      const { data, meta } = await this.fetchServiceDependenciesWithoutStore({
        id,
        params: {
          page: newPage,
          limit: PAGINATION_LIMIT,
          with_flags: true,
          with_state_setting: true,
          define_state: true,
        },
      });

      target.data({
        pending: false,
      });

      this.$set(this.metaByEntityId, id, meta);

      const ids = data.map((item) => {
        let newEntityItem = { entity: item };

        if (this.entitiesById[item._id]) {
          newEntityItem = {
            ...this.entitiesById[item._id],

            entity: {
              ...newEntityItem,
              ...this.entitiesById[item._id].entity,
            },
          };
        }

        this.$set(this.entitiesById, item._id, newEntityItem);

        return item._id;
      });

      const previousDeps = this.entitiesById[id].dependencies ?? [];

      this.$set(this.entitiesById[id], 'dependencies', [
        ...previousDeps,
        ...ids,
      ]);

      this.resetLayout();
    },

    hideDependencies(target) {
      const { entity } = target.data();

      this.$set(this.entitiesById[entity._id], 'dependencies', []);
      this.$delete(this.metaByEntityId, entity._id);

      this.resetLayout();
    },

    /**
     * Method for dependencies fetching for special node
     *
     * @param {string} id
     */
    async toggleDependencies(id) {
      const target = this.$refs.networkGraph.$cy.getElementById(id);
      const { opened, root } = target.data();

      if (!root && opened) {
        this.hideDependencies(target);
      } else {
        await this.showDependencies(target);
      }

      this.runLayout();
    },

    async fetchDependencies(id) {
      const target = this.$refs.networkGraph.$cy.getElementById(id);

      await this.showDependencies(target);

      this.runLayout();
    },

    /**
     * Handler for tap event on whole cytoscape canvas
     *
     * @param {Object} target
     * @param {MouseEvent} originalEvent
     */
    tapHandler({ target, originalEvent }) {
      const { entity, showMore, pending, cycle } = target.data();

      if (cycle || pending) {
        return;
      }

      if (originalEvent.target.classList.contains('v-badge__badge')) {
        const { id } = originalEvent.target.dataset;

        if (id) {
          this.toggleDependencies(id);

          return;
        }
      }

      if (!showMore || !entity) {
        return;
      }

      this.fetchDependencies(entity._id);
    },
  },
};
</script>

<style lang="scss">
.entity-dependencies-by-state-settings {
  position: relative;
  height: 650px;
  width: 100%;
  border-radius: 5px;
  background: white;

  &-network-graph {
    opacity: 0;

    &--ready {
      opacity: 1;
    }
  }

  canvas[data-id='layer0-selectbox'] { // Hide selectbox layer from cytoscape
    display: none;
  }

  .v-badge__badge {
    top: -7px;
    right: -7px;

    * {
      pointer-events: none;
    }
  }
}
</style>
