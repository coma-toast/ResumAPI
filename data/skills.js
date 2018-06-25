'use strict';
var Mockgen = require('./mockgen.js');
/**
 * Operations on /skills
 */
module.exports = {
    /**
     * summary: Retrieves skills.
     * description: 
     * parameters: 
     * produces: application/json, text/json
     * responses: 200
     * operationId: skills_get
     */
    get: {
        200: function (req, res, callback) {
            /**
             * Using mock data generator module.
             * Replace this by actual data for the api.
             */
            Mockgen().responses({
                path: '/skills',
                operation: 'get',
                response: '200'
            }, callback);
        }
    }
};
