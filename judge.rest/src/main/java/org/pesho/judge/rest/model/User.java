package org.pesho.judge.rest.model;

import java.io.Serializable;
import java.sql.Date;
import java.util.List;

import javax.persistence.Basic;
import javax.persistence.Column;
import javax.persistence.Entity;
import javax.persistence.FetchType;
import javax.persistence.GeneratedValue;
import javax.persistence.Id;
import javax.persistence.JoinColumn;
import javax.persistence.ManyToMany;
import javax.persistence.Table;

@Entity
@Table(name = "users")
public class User implements Serializable {

	private static final long serialVersionUID = 1L;

	@Id
	@GeneratedValue
	@Column(name = "id")
	private int id;

	@ManyToMany(fetch=FetchType.EAGER)
	@JoinColumn(name="role_id")
	private List<Role> roles;
	
	@Column(name = "username", unique=true, length = 50)
	private String username;

	@Basic(optional = false)
	@Column(name = "firstname", length = 50)
	private String firstname;

	@Basic(optional = false)
	@Column(name = "lastname", length = 50)
	private String lastname;

	@Basic(optional = false)
	@Column(name = "email", length = 50)
	private String email;

	@Basic(optional = false)
	@Column(name = "active")
	private boolean active;

	@Column(name = "time")
	private Date time;

	public User() {
	}
	
	public User(String username, String firstname, String lastname, String email) {
		this.username = username;
		this.firstname = firstname;
		this.lastname = lastname;
		this.email = email;
	}
	
	public List<Role> getRoles() {
		return roles;
	}
	
	public void setRoles(List<Role> roles) {
		this.roles = roles;
	}

	public int getId() {
		return id;
	}

	public void setId(int id) {
		this.id = id;
	}

	public String getUsername() {
		return username;
	}

	public void setUsername(String username) {
		this.username = username;
	}

	public String getFirstname() {
		return firstname;
	}

	public void setFirstname(String firstname) {
		this.firstname = firstname;
	}

	public String getLastname() {
		return lastname;
	}

	public void setLastname(String lastname) {
		this.lastname = lastname;
	}

	public String getEmail() {
		return email;
	}

	public void setEmail(String email) {
		this.email = email;
	}

	public boolean isActive() {
		return active;
	}

	public void setActive(boolean active) {
		this.active = active;
	}

	public Date getTime() {
		return time;
	}

	public void setTime(Date time) {
		this.time = time;
	}

}