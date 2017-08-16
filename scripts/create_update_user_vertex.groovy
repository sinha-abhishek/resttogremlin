v1=g.V().has('number', XNUM);
if (v1.toList().size() == 0) {
  println('create vertex '+ XNUM );
   v = graph.addVertex('user');
   v.property('number',XNUM);
   v.property('uid',XUID);
} else {
  println('found v');
  v = g.V().has('number', XNUM).next();
  v.property('uid',XUID);
  g.V(v).inE('contact').property('has_app', true)
}
