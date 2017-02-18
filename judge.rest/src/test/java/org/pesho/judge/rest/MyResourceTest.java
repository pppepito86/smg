package org.pesho.judge.rest;

import javax.ws.rs.client.Client;
import javax.ws.rs.client.ClientBuilder;
import javax.ws.rs.client.Entity;
import javax.ws.rs.client.WebTarget;
import javax.ws.rs.core.Response;

import org.glassfish.grizzly.http.server.HttpServer;
import org.junit.After;
import org.junit.Before;
import org.junit.Test;
import org.pesho.judge.rest.model.Role;

import static org.junit.Assert.assertEquals;
import static org.junit.Assert.assertNotSame;

public class MyResourceTest {

    private HttpServer server;
    private WebTarget target;

    @Before
    public void setUp() throws Exception {
        server = Main.startServer();
        Client c = ClientBuilder.newClient();
        target = c.target(Main.BASE_URI);
    }

    @After
    public void tearDown() throws Exception {
        server.stop();
    }

    @Test
    public void testGetIt() {
        String responseMsg = target.path("myresource").request().get(String.class);
        assertEquals("Got it!", responseMsg);
    }
    
    @Test
    public void testCreateRole() {
    	Role role = new Role();
    	role.setRoleName("admin");
    	role.setDescription("admin role");
    	Response response = target.path("roles").request().accept("application/json").post(Entity.json(role));
    	Role responseRole = (Role) response.readEntity(Role.class);
        assertNotSame(responseRole.getId(), 0);
    }
}
