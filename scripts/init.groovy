USER_LABEL = "user";
NUMBER_PROPERTY = "number";
UID_PROPERTY = "uid";
CONTACT_EDGE_LABEL = "contact";
KNOWS_AS_EDGE_PROPERTY = "knw_as";
HAS_APP_PROPERTY = "has_app";
NUMBER_INDEX = "number_index";
UID_INDEX = "uid_index";

changed = false;
mgmt = graph.openManagement();
userLabelCreated = mgmt.getVertexLabel(USER_LABEL) != null;
edgeLabelCreated = mgmt.getEdgeLabel(CONTACT_EDGE_LABEL) != null;
numberIndexCreated = mgmt.getGraphIndex(NUMBER_INDEX) != null;
uidIndexCreated = mgmt.getGraphIndex(UID_INDEX) != null;
mgmt.commit();
if (!userLabelCreated) {
    graph.tx().rollback();
    mgmt = graph.openManagement();
    mgmt.makeVertexLabel(USER_LABEL).make();
    mgmt.commit();
}

if (!edgeLabelCreated) {
    graph.tx().rollback();
    mgmt = graph.openManagement();
    mgmt.makePropertyKey(KNOWS_AS_EDGE_PROPERTY).dataType(String.class).make();
    mgmt.makePropertyKey(HAS_APP_PROPERTY).dataType(Boolean.class).make();
    mgmt.makeEdgeLabel(CONTACT_EDGE_LABEL).make();
    mgmt.commit();
    changed = true;

}


if (!numberIndexCreated) {
    graph.tx().rollback();
    mgmt = graph.openManagement();
    PropertyKey numberPropertyKey = mgmt.makePropertyKey(NUMBER_PROPERTY).dataType(String.class).make();
    mgmt.buildIndex(NUMBER_INDEX, Vertex.class).addKey(numberPropertyKey).buildCompositeIndex();
    changed = true;
    println 'NUMBER_INDEX done';
    mgmt.commit();
}

if (!uidIndexCreated) {
    graph.tx().rollback();
    mgmt = graph.openManagement();
    PropertyKey uidProperty = mgmt.makePropertyKey(UID_PROPERTY).dataType(Integer.class).make();
    mgmt.buildIndex(UID_INDEX, Vertex.class).addKey(uidProperty).buildCompositeIndex();
    changed = true;
    println 'UID_INDEX done';
    mgmt.commit();
}

println 'done'

mgmt.awaitGraphIndexStatus(graph, NUMBER_INDEX).status(SchemaStatus.ENABLED).call();
mgmt.awaitGraphIndexStatus(graph, UID_INDEX).status(SchemaStatus.ENABLED).call();
if (changed) {
    mgmt = graph.openManagement();
    mgmt.updateIndex(mgmt.getGraphIndex(NUMBER_INDEX), SchemaAction.REINDEX).get();
    mgmt.updateIndex(mgmt.getGraphIndex(UID_INDEX), SchemaAction.REINDEX).get();
    mgmt.commit();
    graph.tx().commit();
}
