/**
 * @typedef {Object} AlarmLink
 * @property {string} label
 * @property {string} link
 */

/**
 * @typedef {Object.<string, AlarmLink[] | string[]>} AlarmLinks
 */

/**
 * @typedef {Object} AlarmEvent
 * @property {number} val
 * @property {string} a
 * @property {number} t
 * @property {string} _t
 * @property {string} m
 */

/**
 * @typedef {Object} Correlation
 * @property {number} total
 * @property {Alarm} data
 */

/**
 * @typedef {Object} AlarmRule
 * @property {string} id
 * @property {string} name
 */

/**
 * @typedef {Object} AlarmValue
 * @property {string[]} long_output_history
 * @property {number} last_event_date
 * @property {string} connector_name
 * @property {string} initial_long_output
 * @property {string} output
 * @property {string} initial_output
 * @property {string[]} children
 * @property {string[]} filtered_children
 * @property {Object} extra
 * @property {number} last_update_date
 * @property {number} resolved
 * @property {string} resource
 * @property {number} creation_date
 * @property {string} display_name
 * @property {number} total_state_changes
 * @property {string[]} tags
 * @property {string} long_output
 * @property {string} component
 * @property {string} connector
 * @property {Object} infos_rule_version
 * @property {InfosObject} infos
 * @property {string[]} parents
 * @property {AlarmEvent} ack
 * @property {AlarmEvent} state
 * @property {AlarmEvent} status
 * @property {AlarmEvent} ticket
 * @property {AlarmEvent} snooze
 * @property {AlarmEvent} canceled
 * @property {AlarmEvent} lastComment
 * @property {AlarmEvent[]} [steps]
 */

/**
 * @typedef {Object} Alarm
 * @property {string} _id
 * @property {Entity} entity
 * @property {boolean} metaalarm
 * @property {InfosObject} infos
 * @property {AlarmRule} rule
 * @property {Correlation} consequences
 * @property {Correlation} causes
 * @property {Pbehavior} pbehavior
 * @property {string} links
 * @property {number} t
 * @property {AlarmValue} v
 */