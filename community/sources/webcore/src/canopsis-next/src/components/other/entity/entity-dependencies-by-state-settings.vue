<template>
  <div class="entity-dependencies-by-state-settings">
    <c-zoom-overlay>
      <c-progress-overlay :pending="!ready || pending" />
      <network-graph
        ref="networkGraph"
        :options="options"
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

import { PAGINATION_LIMIT, VUETIFY_ANIMATION_DELAY } from '@/config';
import { ROOT_CAUSE_DIAGRAM_OPTIONS, STATE_SETTING_METHODS, ROOT_CAUSE_DIAGRAM_LAYOUT_OPTIONS } from '@/constants';

import { normalizeTreeOfDependenciesMapEntities } from '@/helpers/entities/map/list';
import { isEntityEventsStateSettings } from '@/helpers/entities/entity/entity';
import { convertSortToRequest } from '@/helpers/entities/shared/query';
import { getButtonHTML, getEntityNodeElementHTML } from '@/helpers/entities/entity/cytoscape';

import { entitiesEntityDependenciesMixin } from '@/mixins/entities/entity-dependencies';

import NetworkGraph from '@/components/common/chart/network-graph.vue';

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
      metaByEntityId: {},
      entitiesById: normalizeTreeOfDependenciesMapEntities([{ entity: this.entity, pinned_entities: [] }]),
    };
  },
  computed: {
    stateSetting() {
      return this.entity?.stateSetting;
    },

    isEventsStateSettings() {
      return isEntityEventsStateSettings(this.entity);
    },

    isInheritedMethod() {
      return this.stateSetting?.method === STATE_SETTING_METHODS.inherited;
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

      if (isEntityEventsStateSettings(entity)) {
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
      const getShowMoreContent = (node) => {
        const { entity } = node;
        const meta = this.metaByEntityId[entity._id] ?? {};

        const fetchedEntities = meta.page * meta.per_page;

        /**
         * TODO: Should be replaced on translation
         */
        return getButtonHTML(`Show more (${fetchedEntities} of ${meta.total_count})`);
      };

      return [
        {
          query: 'node',
          valign: 'center',
          halign: 'center',
          tpl: getEntityNodeElementHTML,
        },
        {
          query: 'node[showMore]',
          valign: 'center',
          halign: 'center',
          tpl: getShowMoreContent,
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

        if (isEntityEventsStateSettings(child)) {
          acc.push(...this.getEventsNodeElementByEntity(child));
        }

        return acc;
      }, []);

      dependenciesNodes.push(
        ...this.getShowMoreElements(entity),
      );

      return dependenciesNodes;
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
        }).run();
      } catch (err) {
        console.warn(err);
      }
    },

    getQuery({ page }) {
      const query = {
        page,
        limit: PAGINATION_LIMIT,
        with_flags: true,
        /**
         * TODO: Api doesn't support multi sort
         */
        ...convertSortToRequest(['last_update_date', 'state']),
      };

      if (this.isInheritedMethod) {
        /**
         * TODO: Api doesn't support pattern but we need it (Case 3)
         */
        query.entity_pattern = JSON.stringify(this.stateSetting.inherited_entity_pattern);
      }

      return query;
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
        params: this.getQuery({ page: newPage }),
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
