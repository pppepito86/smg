package org.pesho.judge.rest.resource;

import java.util.List;

import javax.ws.rs.Consumes;
import javax.ws.rs.GET;
import javax.ws.rs.POST;
import javax.ws.rs.PUT;
import javax.ws.rs.Path;
import javax.ws.rs.PathParam;
import javax.ws.rs.Produces;
import javax.ws.rs.core.MediaType;

import org.pesho.judge.rest.model.Assignment;

@Path("assignments")
public class AssignmentsResource {

    @GET
    @Produces(MediaType.APPLICATION_JSON)
    public List<Assignment> listAssignments() {
        return null;
    }
    
    @POST
    @Consumes(MediaType.APPLICATION_JSON)
    @Produces(MediaType.APPLICATION_JSON)
    public Assignment createAssignment(Assignment assignment) {
        return assignment;
    }
    
    @GET
    @Path("{id}")
    @Produces(MediaType.APPLICATION_JSON)
    public Assignment getAssignment(@PathParam("id") int assignmentId) {
        return null;
    }
    
    @PUT
    @Path("{id}")
    @Produces(MediaType.APPLICATION_JSON)
    public Assignment updateAssignment(@PathParam("id") int assignmentId) {
        return null;
    }
}
