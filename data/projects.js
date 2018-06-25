'use strict';
var Mockgen = require('./mockgen.js');
/**
 * Operations on /projects
 */
module.exports = {
    /**
     * summary: Retrieves projects.
     * description: 
     * parameters: 
     * produces: application/json, text/json
     * responses: 200
     * operationId: projects_get
     */
    get: {
        200: function (req, res, callback) {
            /**
             * Using mock data generator module.
             * Replace this by actual data for the api.
             */
            Mockgen().responses({
                path: '/projects',
                operation: 'get',
                response: '200'
            }, callback);
        }
    }
};
