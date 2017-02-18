package org.pesho.judge.rest.resource;

import java.util.List;

import javax.persistence.EntityManager;
import javax.persistence.EntityManagerFactory;
import javax.persistence.Persistence;
import javax.ws.rs.Consumes;
import javax.ws.rs.GET;
import javax.ws.rs.POST;
import javax.ws.rs.PUT;
import javax.ws.rs.Path;
import javax.ws.rs.PathParam;
import javax.ws.rs.Produces;
import javax.ws.rs.core.MediaType;

import org.pesho.judge.rest.model.Role;

@Path("roles")
public class RolesResource {

    @GET
    @Produces(MediaType.APPLICATION_JSON)
    public List<Role> listRoles() {
        return null;
    }
    
    @POST
    @Consumes(MediaType.APPLICATION_JSON)
    @Produces(MediaType.APPLICATION_JSON)
    public Role createRole(Role role) {
    	EntityManagerFactory emf = Persistence.createEntityManagerFactory("judge");
    	EntityManager em = emf.createEntityManager();
    	em.getTransaction().begin();
    	em.persist(role);
    	em.getTransaction().commit();
    	em.clear();
    	emf.close();
    	
        return role;
    }
    
    @GET
    @Path("{id}")
    @Produces(MediaType.APPLICATION_JSON)
    public Role getRole(@PathParam("id") int roleId) {
        return null;
    }
    
    @PUT
    @Path("{id}")
    @Produces(MediaType.APPLICATION_JSON)
    public Role updateRole(@PathParam("id") int roleId) {
        return null;
    }
}
