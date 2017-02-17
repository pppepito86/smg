package org.pesho.judge.rest.model;

import java.io.Serializable;

import javax.persistence.Column;
import javax.persistence.Entity;
import javax.persistence.GeneratedValue;
import javax.persistence.Id;
import javax.persistence.Table;
import javax.validation.constraints.Pattern;
import javax.validation.constraints.Size;

@Entity
@Table(name = "roles")
public class Role implements Serializable {

	private static final long serialVersionUID = 1L;

	@Id
	@GeneratedValue
	@Column(name = "id")
	private int id;

	@Column(name = "role", unique=true, length = 50)
	@Pattern(message="only small letters", regexp="^[a-z]+$")
	@Size(min = 2, max = 50)
	private String role;

	public Role() {
	}
	
	public Role(String role) {
		this.role = role;
	}

	public String getRole() {
		return role;
	}

}